CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    lucide_icon_name TEXT NOT NULL,
    hue INTEGER NOT NULL CHECK (hue >= 0 AND hue <= 360),
    saturation INTEGER NOT NULL CHECK (saturation >= 0 AND saturation <= 100),
    lightness INTEGER NOT NULL CHECK (lightness >= 0 AND lightness <= 100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
