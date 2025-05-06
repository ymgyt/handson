extern crate alloc;

use core::cmp::max;

use alloc::boxed::Box;

use crate::{
    acpi::AcpiRsdpStruct,
    allocator::ALLOCATOR,
    graphics::{draw_test_pattern, fill_rect, Bitmap as _},
    hpet::{set_global_hpet, Hpet},
    info,
    pci::Pci,
    uefi::{
        exit_from_efi_boot_services, EfiHandle, EfiMemoryType, EfiSystemTable, MemoryMapHolder,
        VramBufferInfo,
    },
    x86::{write_cr3, PageAttr, PAGE_SIZE, PML4},
};
pub fn init_basic_runtime(
    image_handle: EfiHandle,
    efi_system_table: &EfiSystemTable,
) -> MemoryMapHolder {
    let mut memory_map = MemoryMapHolder::new();
    exit_from_efi_boot_services(image_handle, efi_system_table, &mut memory_map);

    ALLOCATOR.init_with_mmap(&memory_map);
    memory_map
}

pub fn init_paing(memory_map: &MemoryMapHolder) {
    use crate::uefi::EfiMemoryType::*;
    let mut table = PML4::new();
    let mut end_of_mem = 0x1_000_000u64;
    for e in memory_map.iter() {
        match e.memory_type {
            CONVENTIONAL_MEMORY | LOADER_CODE | LOADER_DATA => {
                end_of_mem = max(
                    end_of_mem,
                    e.physical_start + e.number_of_pages * (PAGE_SIZE as u64),
                );
            }
            _ => (),
        }
    }
    table
        .create_mapping(0, end_of_mem, 0, PageAttr::ReadWriteKernel)
        .expect("Failed to create intial page mapping");

    // Unmap page 0 to detect null ptr dereference
    table
        .create_mapping(0, 4096, 0, PageAttr::NotPresent)
        .expect("Failed to unmap page 0");
    unsafe {
        write_cr3(Box::into_raw(table));
    }
}

pub fn init_hpet(acpi: &AcpiRsdpStruct) {
    let hpet = acpi.hpet().expect("Failed to get HPET from ACPI");
    let hpet = hpet
        .base_address()
        .expect("Failed to get HPET base address");
    info!("HPET is at {hpet:#p}");
    let hpet = Hpet::new(hpet);
    set_global_hpet(hpet);
}

pub fn init_allocator(memory_map: &MemoryMapHolder) {
    let mut total_memory_pages = 0;
    for e in memory_map.iter() {
        if e.memory_type != EfiMemoryType::CONVENTIONAL_MEMORY {
            continue;
        }
        total_memory_pages += e.number_of_pages;
        info!("{e:?}");
    }
    let total_memory_size_mib = total_memory_pages * 4096 / 1024 / 1024;
    info!("Total: {total_memory_pages} pages = {total_memory_size_mib} MiB");
}

pub fn init_display(vram: &mut VramBufferInfo) {
    let vw = vram.width();
    let vh = vram.height();
    fill_rect(vram, 0x000000, 0, 0, vw, vh).expect("fill_rect failed");
    draw_test_pattern(vram);
}

pub fn init_pci(acpi: &AcpiRsdpStruct) {
    if let Some(mcfg) = acpi.mcfg() {
        for i in 0..mcfg.num_of_entries() {
            if let Some(e) = mcfg.entry(i) {
                info!("{}", e);
            }
        }
        let pci = Pci::new(mcfg);
        pci.probe_devices();
    }
}
