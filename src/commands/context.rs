use std::str::FromStr;

use tetka::games::{games::chess, interface::PositionType};

use crate::core::Searcher;

pub struct Context {
    pub position: chess::Position,
    pub searcher: Searcher,
}

impl Default for Context {
    fn default() -> Self {
        let position = chess::Position::from_str(chess::Position::STARTPOS).unwrap();
        Context {
            position,
            searcher: Searcher::new(),
        }
    }
}
