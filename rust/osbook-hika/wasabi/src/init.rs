extern crate alloc;

use core::cmp::max;

use alloc::boxed::Box;

use crate::{
    allocator::ALLOCATOR,
    uefi::{exit_from_efi_boot_services, EfiHandle, EfiSystemTable, MemoryMapHolder},
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
    unsafe {
        write_cr3(Box::into_raw(table));
    }
}
