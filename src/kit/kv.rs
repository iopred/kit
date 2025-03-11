#[derive(Clone, Debug)]
pub enum KitValue {
    String(String),
    Byte(u8),
    Boolean(bool),
}

impl PartialEq for KitValue {
    fn eq(&self, other: &Self) -> bool {
        self.contains(other)
    }
}

// [🌌] Branch when the current universe doesn't match.
// [🔍] Identify the cause of the disturbance.
// [⛔] Stop the simulation when a disturbance is detected.
// [🔼] Move up.
impl KitValue {
    fn contains(&self, other: &KitValue) -> bool {
        match (self, other) {
            (KitValue::String(a), KitValue::String(b)) => a.contains(b),
            (KitValue::String(a), KitValue::Byte(b)) => a.contains(&b.to_string()),
            (KitValue::Byte(a), KitValue::String(b)) => a.to_string().contains(b),
            (KitValue::Byte(a), KitValue::Byte(b)) => a.to_string().contains(&b.to_string()),
            _ => self == other,
        }
    }
}