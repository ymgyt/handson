#![feature(prelude_import)]
#![no_std]
#![no_main]
#[prelude_import]
use core::prelude::rust_2021::*;
#[macro_use]
extern crate core;
extern crate compiler_builtins as _;
use aya_ebpf::{macros::uprobe, programs::ProbeContext};
use aya_log_ebpf::info;
#[no_mangle]
#[link_section = "uprobe"]
pub fn uprobeho(ctx: *mut ::core::ffi::c_void) -> u32 {
    let _ = uprobeho(::aya_ebpf::programs::ProbeContext::new(ctx));
    return 0;
    pub fn uprobeho(ctx: ProbeContext) -> u32 {
        match try_uprobeho(ctx) {
            Ok(ret) => ret,
            Err(ret) => ret,
        }
    }
}
fn try_uprobeho(ctx: ProbeContext) -> Result<u32, u32> {
    match unsafe { &mut ::aya_log_ebpf::AYA_LOG_BUF }
        .get_ptr_mut(0)
        .and_then(|ptr| unsafe { ptr.as_mut() })
    {
        None => {}
        Some(::aya_log_ebpf::LogBuf { buf }) => {
            let _: Option<()> = (|| {
                let size = ::aya_log_ebpf::write_record_header(
                    buf,
                    "uprobeho",
                    ::aya_log_ebpf::macro_support::Level::Info,
                    "uprobeho",
                    "uprobeho-ebpf/src/main.rs",
                    16u32,
                    1usize,
                )?;
                let mut size = size.get();
                let slice = buf.get_mut(size..)?;
                let len = ::aya_log_ebpf::WriteToBuf::write(
                    "function uprobed_function called by /proc/self/exe",
                    slice,
                )?;
                size += len.get();
                let record = buf.get(..size)?;
                unsafe { &mut ::aya_log_ebpf::AYA_LOGS }.output(&ctx, record, 0);
                Some(())
            })();
        }
    };
    let arg: u32 = ctx.arg(0).unwrap_or(9999);
    match unsafe { &mut ::aya_log_ebpf::AYA_LOG_BUF }
        .get_ptr_mut(0)
        .and_then(|ptr| unsafe { ptr.as_mut() })
    {
        None => {}
        Some(::aya_log_ebpf::LogBuf { buf }) => {
            let _: Option<()> = (|| {
                let size = ::aya_log_ebpf::write_record_header(
                    buf,
                    "uprobeho",
                    ::aya_log_ebpf::macro_support::Level::Info,
                    "uprobeho",
                    "uprobeho-ebpf/src/main.rs",
                    18u32,
                    3usize,
                )?;
                let mut size = size.get();
                let slice = buf.get_mut(size..)?;
                let len = ::aya_log_ebpf::WriteToBuf::write("arg: ", slice)?;
                size += len.get();
                let slice = buf.get_mut(size..)?;
                let len = ::aya_log_ebpf::WriteToBuf::write(
                    ::aya_log_ebpf::macro_support::DisplayHint::Default,
                    slice,
                )?;
                size += len.get();
                let slice = buf.get_mut(size..)?;
                let len = ::aya_log_ebpf::WriteToBuf::write(
                    {
                        let tmp = arg;
                        let _: &dyn ::aya_log_ebpf::macro_support::DefaultFormatter = &tmp;
                        tmp
                    },
                    slice,
                )?;
                size += len.get();
                let record = buf.get(..size)?;
                unsafe { &mut ::aya_log_ebpf::AYA_LOGS }.output(&ctx, record, 0);
                Some(())
            })();
        }
    };
    Ok(0)
}
#[cfg(not(test))]
#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    loop {}
}
