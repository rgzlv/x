//! Functions for file operations related to the
//! [XDG Base Directory specification] and the environments.
//!
//!
//! [XDG Base Directory specification]: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html

use std::env;
use std::io;
use std::path;
use std::result;

use thiserror::Error;
use toml::ser;

use crate::debug::DebugPanic;
use crate::table;

type Result<T> = result::Result<T, Err>;

#[non_exhaustive]
#[derive(Debug, Error)]
pub enum Err {
	#[error(transparent)]
	Table(#[from] Box<table::Err>),
	#[error(transparent)]
	IO(#[from] io::Error),
	#[error(transparent)]
	VarErr(#[from] env::VarError),
	#[error(transparent)]
	TOMLSerializeErr(#[from] ser::Error),

	#[error("Couldn't convert path to string.")]
	PathToStr,
	#[error("Binary '{0}' not in PATH. If it does exist, check permissions.")]
	NoBinInPath(path::PathBuf),
	#[error("binary '{0}' doesn't exist on host")]
	NoExistsBin(path::PathBuf),
	#[error("Directory '{0}' already exists")]
	DirExists(path::PathBuf),
}

pub fn lookup_bin(bin: &path::Path) -> Result<path::PathBuf> {
	let path = env::var("PATH").dp()?;
	let path = path.split(':');

	for path_elem in path {
		let path_elem = path::PathBuf::from(path_elem);
		let bin_in_path_elem = path_elem.join(bin);
		if let Ok(exists) = bin_in_path_elem.try_exists() {
			if exists {
				return Ok(bin_in_path_elem);
			}
		}
		// Probably don't want to error if there's no permission to access
		// a path element. `Err::NoBinInPath` also says to check permissions
		// to PATH if `bin` couldn't be found.
	}

	Err(Err::NoBinInPath(bin.to_owned()))
}

pub fn bin_try_exists(bin: &path::Path) -> Result<()> {
	match bin.try_exists() {
		Ok(exists) => {
			if !exists {
				return Err(Err::NoExistsBin(bin.to_owned())).dp();
			}
		}
		Err(err) => {
			Err(err).dp()?;
		}
	}
	Ok(())
}

pub fn bin_get_abs(bin: &path::Path) -> Result<path::PathBuf> {
	if bin.is_absolute() {
		Ok(bin.to_owned())
	} else {
		lookup_bin(bin)
	}
}
