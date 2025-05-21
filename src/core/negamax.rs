use tetka::games::{games::chess, interface::PositionType};

use super::Searcher;

pub fn negamax(context: &mut Searcher, plys: u8, depth: u8, alpha: u16, beta: u16) -> u16 {
    if depth == 0 {
        return 0;
    }
    context.register_depth(plys);

    if context.should_stop() {
        return 0;
    }

    let moves = context.position.generate_moves::<false, true, true>();
    for mv in moves {
        let score = context
            .position
            .insert(chess::Square::A1, chess::ColoredPiece::BlackBishop);
    }

    0
}
