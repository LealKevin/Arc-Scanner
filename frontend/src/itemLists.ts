export type WorkshopRequirement = {
  workshop: string;
  level: number;
};

export const WORKSHOP_ITEMS: Record<string, WorkshopRequirement[]> = {
  // Scrappy (L2-L4)
  "dog-collar": [{ workshop: "Scrappy", level: 2 }],
  "lemon": [{ workshop: "Scrappy", level: 2 }],
  "apricot": [
    { workshop: "Scrappy", level: 2 },
    { workshop: "Scrappy", level: 4 },
  ],
  "prickly-pear": [{ workshop: "Scrappy", level: 3 }],
  "olives": [{ workshop: "Scrappy", level: 3 }],
  "cat-bed": [{ workshop: "Scrappy", level: 3 }],
  "mushroom": [{ workshop: "Scrappy", level: 4 }],
  "very-comfortable-pillow": [{ workshop: "Scrappy", level: 4 }],

  // Gunsmith (L2-L3)
  "rusted-tools": [{ workshop: "Gunsmith", level: 2 }],
  "mechanical-components": [{ workshop: "Gunsmith", level: 2 }],
  "wasp-driver": [{ workshop: "Gunsmith", level: 2 }],
  "rusted-gear": [{ workshop: "Gunsmith", level: 3 }],
  "advanced-mechanical-components": [{ workshop: "Gunsmith", level: 3 }],
  "sentinel-part": [{ workshop: "Gunsmith", level: 3 }],

  // Medical Lab (L2-L3)
  "cracked-bioscanner": [{ workshop: "Medical Lab", level: 2 }],
  "durable-cloth": [{ workshop: "Medical Lab", level: 2 }],
  "tick-pod": [{ workshop: "Medical Lab", level: 2 }],
  "rusted-shut-medical-kit": [{ workshop: "Medical Lab", level: 3 }],
  "antiseptic": [{ workshop: "Medical Lab", level: 3 }],
  "surveyor-vault": [{ workshop: "Medical Lab", level: 3 }],

  // Explosives Station (L2-L3)
  "synthesized-fuel": [{ workshop: "Explosives", level: 2 }],
  "crude-explosives": [{ workshop: "Explosives", level: 2 }],
  "pop-trigger": [{ workshop: "Explosives", level: 2 }],
  "laboratory-reagents": [{ workshop: "Explosives", level: 3 }],
  "explosive-compound": [{ workshop: "Explosives", level: 3 }],
  "rocketeer-part": [{ workshop: "Explosives", level: 3 }],

  // Gear Bench (L2-L3)
  "power-cable": [{ workshop: "Gear Bench", level: 2 }],
  "hornet-driver": [{ workshop: "Gear Bench", level: 2 }],
  "electrical-components": [
    { workshop: "Gear Bench", level: 2 },
    { workshop: "Utility", level: 2 },
  ],
  "industrial-battery": [{ workshop: "Gear Bench", level: 3 }],
  "advanced-electrical-components": [
    { workshop: "Gear Bench", level: 3 },
    { workshop: "Utility", level: 3 },
  ],
  "bastion-part": [{ workshop: "Gear Bench", level: 3 }],

  // Refiner (L2-L3)
  "toaster": [{ workshop: "Refiner", level: 2 }],
  "arc-motion-core": [{ workshop: "Refiner", level: 2 }],
  "fireball-burner": [{ workshop: "Refiner", level: 2 }],
  "motor": [{ workshop: "Refiner", level: 3 }],
  "arc-circuitry": [{ workshop: "Refiner", level: 3 }],
  "bombardier-cell": [{ workshop: "Refiner", level: 3 }],

  // Utility Station (L2-L3)
  "damaged-heat-sink": [{ workshop: "Utility", level: 2 }],
  "snitch-scanner": [{ workshop: "Utility", level: 2 }],
  "fried-motherboard": [{ workshop: "Utility", level: 3 }],
  "leaper-pulse-unit": [{ workshop: "Utility", level: 3 }],
};

// Quest items (from cheat sheet "Keep for Quests")
export const QUEST_ITEM_IDS = [
  "leaper-pulse-unit",
  "power-rod",
  "rocketeer-part",
  "surveyor-vault",
  "antiseptic",
  "hornet-driver",
  "syringe",
  "wasp-driver",
  "water-pump",
  "snitch-scanner",
];

// Project items (from cheat sheet "Keep for Projects")
export const PROJECT_ITEM_IDS = [
  "magnetic-accelerator",
  "exodus-modules",
  "advanced-electrical-components",
  "humidifier",
  "sensors-recipe",
  "cooling-fan",
  "battery",
  "light-bulb",
  "electrical-components",
  "wires-recipe",
  "durable-cloth",
  "spring",
  "arc-alloy",
  "rubber-parts-recipe",
  "metal-parts",
];

// Build KEEP_ITEM_IDS from quest + project + workshop items
export const KEEP_ITEM_IDS: Set<string> = new Set([
  ...QUEST_ITEM_IDS,
  ...PROJECT_ITEM_IDS,
  ...Object.keys(WORKSHOP_ITEMS),
]);

// Safe to Recycle items (from cheat sheet)
export const RECYCLE_ITEM_IDS: Set<string> = new Set([
  "alarm-clock",
  "arc-coolant",
  "arc-flex-rubber",
  "arc-performance-steel",
  "arc-synthetic-resin",
  "arc-thermo-lining",
  "bicycle-pump",
  "broken-flashlight",
  "broken-guidance-system",
  "broken-handcuffs",
  "broken-handheld-radio",
  "broken-taser",
  "burned-arc-circuitry",
  "camera-lens",
  "candle-holder",
  "coolant",
  "cooling-coil",
  "crumpled-plastic-bottle",
  "damaged-arc-motion-core",
  "damaged-arc-powercell",
  "deflated-football",
  "diving-googles",
  "dried-out-arc-resin",
  "expired-respirator",
  "recorder",
  "frying-pan",
  "garlic-press",
  "headphones",
  "ice-cream-scooper",
  "household-cleaner",
  "impure-arc-coolant",
  "industrial-charger",
  "industrial-magnet",
  "metal-brackets",
  "number-plate",
  "polluted-air-filter",
  "portable-television",
  "power-bank",
  "projector",
  "radio",
  "remote-control",
  "ripped-safety-vest",
  "ruined-accordion",
  "ruined-baton",
  "ruined-handcuffs",
  "ruined-parachute",
  "ruined-riot-shield",
  "ruined-tactical-vest",
  "rusted-bolts",
  "rusty-arc-steel",
  "spotter-relay",
  "spring-cushion",
  "tattered-arc-lining",
  "tattered-clothes",
  "thermostat",
  "torn-blanket",
  "turbo-pump",
  "water-filter",
]);

export type ItemCategory = "keep" | "recycle" | "unknown";

export function getItemCategory(itemId: string): ItemCategory {
  if (KEEP_ITEM_IDS.has(itemId)) return "keep";
  if (RECYCLE_ITEM_IDS.has(itemId)) return "recycle";
  return "unknown";
}

export function getWorkshopRequirements(itemId: string): WorkshopRequirement[] {
  return WORKSHOP_ITEMS[itemId] || [];
}

export function getItemTags(itemId: string): string[] {
  const tags: string[] = [];
  if (QUEST_ITEM_IDS.includes(itemId)) tags.push("Quest");
  if (PROJECT_ITEM_IDS.includes(itemId)) tags.push("Project");
  return tags;
}
