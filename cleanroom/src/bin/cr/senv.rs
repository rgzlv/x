//! Structs and methods for storing and operating on files related to the shell
//! environments.

use std::fs;
use std::io;
use std::path;
use std::result;

use thiserror::Error;
use toml::ser;

use crate::debug::{dbgfmt, DebugPanic};
use crate::table;

type Result<T> = result::Result<T, Err>;

#[non_exhaustive]
#[derive(Debug, Error)]
pub enum Err {
	#[error(transparent)]
	IO(#[from] io::Error),
	#[error(transparent)]
	TOMLSerialize(#[from] ser::Error),

	#[error("Directory '{1}' doesn't exist for environment '{0}'")]
	MissingDir(String, path::PathBuf),
	#[error("File '{1}' doesn't exist for environment '{0}'")]
	MissingFile(String, path::PathBuf),
}

use std::cmp::{Eq, Ord, PartialEq, PartialOrd};

#[non_exhaustive]
#[derive(Eq, Ord, PartialEq, PartialOrd, Debug)]
pub struct Senv {
	pub name: String,
	pub files: Files,
}

#[non_exhaustive]
#[derive(Eq, Ord, PartialEq, PartialOrd, Debug)]
pub struct Files {
	pub cfg_dir: path::PathBuf,
	pub cfg_file: path::PathBuf,
	pub data_dir: path::PathBuf,
	pub bin_dir: path::PathBuf,
}

impl Senv {
	pub fn new_xdg(name: &str, dirs: &xdg::BaseDirectories) -> Result<Self> {
		let name = String::from(name);
		let cfg_dir = dirs.get_config_home().join(&name);
		let cfg_file = cfg_dir.join("config.toml");
		let data_dir = dirs.get_data_home().join(&name);
		let bin_dir = data_dir.join("bin");

		Ok(Self {
			name,
			files: Files {
				cfg_dir,
				cfg_file,
				data_dir,
				bin_dir,
			},
		})
	}

	pub fn create_xdg(self) -> Result<Self> {
		fs::create_dir_all(&self.files.cfg_dir).dp()?;
		fs::create_dir_all(&self.files.data_dir).dp()?;
		fs::create_dir_all(&self.files.bin_dir).dp()?;
		fs::File::create_new(&self.files.cfg_file).dp()?;
		fs::write(
			&self.files.cfg_file,
			toml::to_string_pretty(&table::Root::new()).dp()?,
		)
		.dp()?;
		Ok(self)
	}

	pub fn create_new_xdg(
		name: &str,
		dirs: &xdg::BaseDirectories,
	) -> Result<Self> {
		println!("{:?}", dirs);
		Self::new_xdg(name, dirs).dp()?.create_xdg().dp()
	}

	pub fn rm(self) -> Result<()> {
		fs::remove_dir_all(&self.files.cfg_dir).dp()?;
		fs::remove_dir_all(&self.files.data_dir).dp()?;

		Ok(())
	}

	pub fn is_valid(&self) -> Result<()> {
		if !self.files.cfg_dir.try_exists().dp()? {
			return Err(Err::MissingDir(
				self.name.clone(),
				self.files.cfg_dir.clone(),
			))
			.dp();
		}

		if !self.files.cfg_file.try_exists().dp()? {
			return Err(Err::MissingFile(
				self.name.clone(),
				self.files.cfg_file.clone(),
			))
			.dp();
		}

		if !self.files.data_dir.try_exists().dp()? {
			return Err(Err::MissingDir(
				self.name.clone(),
				self.files.data_dir.clone(),
			))
			.dp();
		}

		Ok(())
	}

	#[allow(
		clippy::missing_panics_doc,
		clippy::unwrap_in_result,
		clippy::unwrap_used
	)]
	pub fn get_vec(dirs: &xdg::BaseDirectories) -> Result<Vec<Self>> {
		let mut shell_envs: Vec<Self> = Vec::new();

		let files = fs::read_dir(dirs.get_config_home()).dp()?;
		for file in files {
			if let Err(err) = file {
				dbgfmt!("{}:{} {}", file!(), line!(), err);
				continue;
			}
			let file = file.unwrap();

			let meta = file.metadata();
			if let Err(err) = meta {
				dbgfmt!("{}:{} {}", file!(), line!(), err);
				continue;
			}
			let meta = meta.unwrap();

			if !meta.is_dir() {
				continue;
			}

			let file_name = file.file_name();
			let file_name = file_name.to_str();
			if file_name.is_none() {
				dbgfmt!(
					"{}:{} {}",
					file!(),
					line!(),
					"Couldn't convert file name to `&str`"
				);
				continue;
			}
			let file_name = file_name.unwrap();

			let shell_env = Self::new_xdg(file_name, dirs);
			if shell_env.is_err() {
				continue;
			}
			let shell_env = shell_env.unwrap();

			if shell_env.is_valid().is_err() {
				continue;
			}

			shell_envs.push(shell_env);
		}
		Ok(shell_envs)
	}
}
