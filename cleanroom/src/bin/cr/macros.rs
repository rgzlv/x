/// Like `vec!` but for `std::path::PathBuf`.
macro_rules! pathbuf {
	($($exs:expr),+) => {
		vec![$(std::path::PathBuf::from($exs)),+]
	};
}
pub(crate) use pathbuf;
