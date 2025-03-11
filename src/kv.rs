pub mod kit;

#[derive(Clone, Debug)]
pub enum KitValue {
    String(String),
    Number(f64),
    Boolean(bool),
}

// [🌌] Branch when the current universe doesn't match.
// [🔍] Identify the cause of the disturbance.
// [⛔] Stop the simulation when a disturbance is detected.
// [🔼] Move up.
impl PartialEq for KitValue {
    fn eq(&self, other: &Self) -> bool {
        self.contains(other)
    }
}

impl KitValue {
    fn contains(&self, other: &KitValue) -> bool {
        match (self, other) {
            (KitValue::String(a), KitValue::String(b)) => a.contains(b),
            (KitValue::String(a), KitValue::Number(b)) => a.contains(&b.to_string()),
            (KitValue::Number(a), KitValue::String(b)) => a.to_string().contains(b),
            (KitValue::Number(a), KitValue::Number(b)) => a.to_string().contains(&b.to_string()),
            _ => false,
        }
    }
}