#![no_std]
#![no_main]

use core::fmt::Write as _;
use core::panic::PanicInfo;

use wasabi::{
    error,
    graphics::{draw_test_pattern, fill_rect, Bitmap as _},
    info,
    init::{init_basic_runtime, init_paing},
    print::hexdump,
    println,
    qemu::{exit_qemu, QemuExitCode},
    uefi::{
        init_vram, locate_loaded_image_protocol, EfiHandle, EfiMemoryType, EfiSystemTable,
        VramTextWriter,
    },
    warn,
    x86::hlt,
    x86::{init_exceptions, trigger_debug_interrupt},
};

#[no_mangle]
fn efi_main(image_handle: EfiHandle, efi_system_table: &EfiSystemTable) {
    println!("Booting WasabiOS...");
    println!("image_handle: {:#018X}", image_handle);
    println!("efi_system_table: {:#p}", efi_system_table);
    let loaded_image_protocol = locate_loaded_image_protocol(image_handle, efi_system_table)
        .expect("Failed to get load image protocol");
    println!("image_base: {:018X}", loaded_image_protocol.image_base);
    println!("image_size: {:018X}", loaded_image_protocol.image_size);
    info!("info");
    warn!("warn");
    error!("error");
    hexdump(efi_system_table);
    let mut vram = init_vram(efi_system_table).expect("init_vram failed");
    let vw = vram.width();
    let vh = vram.height();
    fill_rect(&mut vram, 0x000000, 0, 0, vw, vh).expect("fill_rect failed");
    draw_test_pattern(&mut vram);
    let mut w = VramTextWriter::new(&mut vram);
    let memory_map = init_basic_runtime(image_handle, efi_system_table);
    let mut total_memory_pages = 0;
    for e in memory_map.iter() {
        if e.memory_type != EfiMemoryType::CONVENTIONAL_MEMORY {
            continue;
        }
        total_memory_pages += e.number_of_pages;
        writeln!(w, "{e:?}").unwrap();
    }
    let total_memory_size_mib = total_memory_pages * 4096 / 1024 / 1024;
    writeln!(
        w,
        "Total: {total_memory_pages} pages = {total_memory_size_mib} MiB"
    )
    .unwrap();
    writeln!(w, "Hello Non-UEFI world!").unwrap();
    let cr3 = wasabi::x86::read_cr3();
    println!("cr3 = {cr3:#p}");
    let t = Some(unsafe { &*cr3 });
    // println!("{t:?}");
    let t = t.and_then(|t| t.next_level(0));
    // println!("{t:?}");
    let t = t.and_then(|t| t.next_level(0));
    // println!("{t:?}");
    let _t = t.and_then(|t| t.next_level(0));
    // println!("{t:?}");

    let (_gdt, _idt) = init_exceptions();
    info!("Exception initialized!");
    trigger_debug_interrupt();
    info!("Execution continued");
    init_paing(&memory_map);
    info!("Now we are using our own page tables!");
    loop {
        hlt()
    }
}

#[panic_handler]
fn panic(info: &PanicInfo) -> ! {
    error!("PANIC: {info:?}");
    exit_qemu(QemuExitCode::Fail);
}
