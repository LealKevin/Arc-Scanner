import type { Item } from "../types";

type Props = {
  item: Item;
  className?: string;
};

export function ItemCard({ item, className }: Props) {
  return (
    <div className={`item-card ${className ?? ""}`}>
      <img className="item-icon" src={item.icon} alt={item.name} />
      <span className="item-value">$ {item.value}</span>
    </div>
  );
}
