use aya::programs::UProbe;
#[rustfmt::skip]
use log::{info, warn};
use tokio::signal;

#[no_mangle]
#[inline(never)]
pub extern "C" fn uprobed_function(_val: u32) {}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // let opt = Opt::parse();

    env_logger::init();

    // This will include your eBPF object file as raw bytes at compile-time and load it at
    // runtime. This approach is recommended for most real-world use cases. If you would
    // like to specify the eBPF program at runtime rather than at compile-time, you can
    // reach for `Bpf::load_file` instead.
    let bpf_obj_path = concat!(env!("OUT_DIR"), "/uprobeho");
    println!("pbf obj path: {bpf_obj_path}");
    let mut ebpf = aya::Ebpf::load(aya::include_bytes_aligned!(concat!(
        env!("OUT_DIR"),
        "/uprobeho"
    )))?;
    if let Err(e) = aya_log::EbpfLogger::init(&mut ebpf) {
        // This can happen if you remove all log statements from your eBPF program.
        warn!("failed to initialize eBPF logger: {}", e);
    }

    let program: &mut UProbe = ebpf.program_mut("uprobeho").unwrap().try_into()?;
    program.load()?;
    program.attach(Some("uprobed_function"), 0, "/proc/self/exe", None)?;

    info!("Attached!");

    uprobed_function(111);
    uprobed_function(222);

    let ctrl_c = signal::ctrl_c();
    println!("Waiting for Ctrl-C...");
    ctrl_c.await?;
    println!("Exiting...");

    Ok(())
}
