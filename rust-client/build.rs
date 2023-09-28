fn main() {
    println!("cargo:rerun-if-changed=proto/user/user.proto");
    tonic_build::configure()
        .build_server(false)
        .compile(&["proto/user/user.proto"], &["proto"])
        .unwrap_or_else(|e| panic!("Failed to compile protos {:?}", e));
}
