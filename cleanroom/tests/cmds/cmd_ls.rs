use std::path;
use std::process;

use crate::BIN_NAME;

#[test]
fn empty() {
	let xdg_root = path::PathBuf::from(env!("CARGO_TARGET_TMPDIR"));
	let cfg_home = xdg_root.join("config");
	let data_home = xdg_root.join("local").join("share");

	let output = process::Command::new(BIN_NAME)
		.env_clear()
		.env("XDG_CONFIG_HOME", cfg_home.to_str().unwrap())
		.env("XDG_DATA_HOME", data_home.to_str().unwrap())
		.args(["ls"])
		.output()
		.unwrap();
	let output = std::str::from_utf8(&output.stdout).unwrap();

	assert_eq!(output, "");
}

#[test]
fn two() {
	let xdg_root = path::PathBuf::from(env!("CARGO_TARGET_TMPDIR"));
	let cfg_home = xdg_root.join("config");
	let data_home = xdg_root.join("local").join("share");

	let output = process::Command::new(BIN_NAME)
		.env_clear()
		.env("XDG_CONFIG_HOME", cfg_home.to_str().unwrap())
		.env("XDG_DATA_HOME", data_home.to_str().unwrap())
		.args(["ls"])
		.output()
		.unwrap();
	let output = std::str::from_utf8(&output.stdout).unwrap();

	assert_eq!(output, "");
}
