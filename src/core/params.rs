use derive_new::new;

#[derive(Clone)]
pub struct Params {}

impl Params {
    pub fn new() -> Params {
        Params {}
    }
}

#[derive(Clone, new)]
pub struct Param {
    val: f64,

    #[allow(unused)]
    min: f64,
    #[allow(unused)]
    max: f64,
}

impl Param {
    #[allow(unused)]
    pub fn set(&mut self, val: f64) {
        self.val = val;
    }
}
