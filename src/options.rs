use tetka::uxi::Parameter;

pub fn hash() -> Parameter {
    Parameter::Spin(16, 1, 33554432)
}

pub fn threads() -> Parameter {
    Parameter::Spin(1, 1, 1024)
}
