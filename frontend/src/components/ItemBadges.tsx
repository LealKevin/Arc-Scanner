import {
  getItemCategory,
  getWorkshopRequirements,
  getItemTags,
} from "../itemLists";

type Props = {
  itemId: string;
  className?: string;
};

export function ItemBadges({ itemId, className }: Props) {
  const category = getItemCategory(itemId);
  const workshops = getWorkshopRequirements(itemId);
  const tags = getItemTags(itemId);

  if (category === "unknown" && workshops.length === 0 && tags.length === 0) {
    return null;
  }

  return (
    <div className={`item-badges ${className ?? ""}`}>
      {category !== "unknown" && (
        <div
          className={`badge ${category === "keep" ? "unsafeToRecycle" : "safeToRecycle"}`}
        >
          <span>{category === "keep" ? "Keep" : "Recycle"}</span>
        </div>
      )}
      {tags.map((tag, index) => (
        <div key={`tag-${index}`} className={`badge ${tag.toLowerCase()}`}>
          <span>{tag}</span>
        </div>
      ))}
      {workshops.map((req, index) => (
        <div key={`ws-${index}`} className="badge workshop">
          <span>
            {req.workshop} L{req.level}
          </span>
        </div>
      ))}
    </div>
  );
}
