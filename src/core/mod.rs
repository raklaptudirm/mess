use std::time;

use rand::seq::IndexedRandom;
use tetka::games::{
    games::chess,
    interface::{MoveType, PositionType},
};

mod params;
pub use self::params::*;

#[derive(Clone)]
pub struct Searcher {
    limits: Limits,
    params: Params,

    start: time::Instant,

    seldepth: usize,
}

#[derive(Debug, Default, Clone)]
pub struct Limits {
    pub maxdepth: Option<usize>,
    pub maxnodes: Option<usize>,
    pub movetime: Option<u128>,

    #[allow(unused)]
    pub movestogo: Option<usize>,
}

impl Searcher {
    pub fn new() -> Searcher {
        Searcher {
            limits: Default::default(),
            params: Params::new(),
            start: time::Instant::now(),
            seldepth: 0,
        }
    }

    pub fn search(
        &mut self,
        position: &chess::Position,
        limits: Limits,
        total_nodes: &mut u64,
    ) -> chess::Move {
        self.start = time::Instant::now();
        self.seldepth = 0;

        let mut move_vec = Vec::<chess::Move>::new();
        let moves = position.generate_moves::<false, true, true>();
        for mv in moves {
            move_vec.push(mv);
        }

        match move_vec.choose(&mut rand::rng()) {
            Some(mv) => *mv,
            None => chess::Move::NULL,
        }
    }
}
