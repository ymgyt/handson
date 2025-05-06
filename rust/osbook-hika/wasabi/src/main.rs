#![no_std]
#![no_main]

use core::panic::PanicInfo;
use core::time::Duration;

use wasabi::{
    error,
    executor::{Executor, Task, TimeoutFuture},
    hpet::global_timestamp,
    info,
    init::{init_allocator, init_basic_runtime, init_display, init_hpet, init_paing},
    print::{hexdump, set_global_vram},
    println,
    qemu::{exit_qemu, QemuExitCode},
    serial::SerialPort,
    uefi::{init_vram, locate_loaded_image_protocol, EfiHandle, EfiSystemTable},
    x86::init_exceptions,
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
    hexdump(efi_system_table);

    let mut vram = init_vram(efi_system_table).expect("init_vram failed");
    init_display(&mut vram);
    set_global_vram(vram);

    let acpi = efi_system_table.acpi_table().expect("ACPI table not found");
    let memory_map = init_basic_runtime(image_handle, efi_system_table);
    info!("Hello Non-UEFI world!");
    init_allocator(&memory_map);

    let (_gdt, _idt) = init_exceptions();
    init_paing(&memory_map);
    init_hpet(acpi);

    let t0 = global_timestamp();
    let serial_task = Task::new(async {
        let sp = SerialPort::default();
        if let Err(e) = sp.loopback_test() {
            error!("{e:?}");
            return Err("serial: llopback test failed");
        }
        info!("Started to monitor serial port");
        loop {
            if let Some(v) = sp.try_read() {
                let c = char::from_u32(v as u32);
                info!("serial input: {v:#04X} = {c:?}");
            }
            TimeoutFuture::new(Duration::from_millis(20)).await;
        }
    });
    let task1 = Task::new(async move {
        for i in 100..=103 {
            info!("{i} hpet.main_counter = {:?}", global_timestamp() - t0);
            TimeoutFuture::new(Duration::from_secs(1)).await;
        }
        Ok(())
    });
    let task2 = Task::new(async move {
        for i in 200..=203 {
            info!("{i} hpet.main_counter = {:?}", global_timestamp() - t0);
            TimeoutFuture::new(Duration::from_secs(1)).await;
        }
        Ok(())
    });
    let mut executor = Executor::new();
    executor.enqueue(task1);
    executor.enqueue(task2);
    executor.enqueue(serial_task);
    Executor::run(executor)
}

#[panic_handler]
fn panic(info: &PanicInfo) -> ! {
    error!("PANIC: {info:?}");
    exit_qemu(QemuExitCode::Fail);
}
