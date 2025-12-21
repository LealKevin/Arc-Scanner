import { useEffect, useState, useRef } from "react";
import "./App.css";
import { getItemCategory, getWorkshopRequirements, getItemTags } from "./itemLists";
import {
  EventsOn,
  WindowSetSize,
  WindowSetPosition,
  ScreenGetAll,
} from "../wailsjs/runtime/runtime";

const WINDOW_WIDTH_VISIBLE = 100;
const WINDOW_HEIGHT_VISIBLE = 220;
const WINDOW_WIDTH_HIDDEN = 1;
const WINDOW_HEIGHT_HIDDEN = 1;

type Item = {
  id: string;
  name: string;
  value: number;
  icon: string;
};

function App() {
  const [item, setItem] = useState<Item>();
  const [isScanning, setIsScanning] = useState(false);
  const [isScanningFailed, setIsScanningFailed] = useState(false);
  const [isVisible, setIsVisible] = useState(true);
  const [showItem, setShowItem] = useState(false);

  const fadeTimeoutRef = useRef<number | null>(null);
  const clearTimeoutRef = useRef<number | null>(null);

  const updateWindowSize = async (visible: boolean) => {
    const screens = await ScreenGetAll();
    const currentScreen = screens.find((s) => s.isCurrent) || screens[0];

    if (visible) {
      const x = currentScreen.width - WINDOW_WIDTH_VISIBLE;
      WindowSetSize(WINDOW_WIDTH_VISIBLE, WINDOW_HEIGHT_VISIBLE);
      WindowSetPosition(x, 0);
    } else {
      const x = currentScreen.width - WINDOW_WIDTH_HIDDEN;
      WindowSetSize(WINDOW_WIDTH_HIDDEN, WINDOW_HEIGHT_HIDDEN);
      WindowSetPosition(x, 0);
    }
  };

  useEffect(() => {
    // Start with collapsed window
    updateWindowSize(false);

    EventsOn("item-found", (data) => {
      if (fadeTimeoutRef.current) window.clearTimeout(fadeTimeoutRef.current);
      if (clearTimeoutRef.current) window.clearTimeout(clearTimeoutRef.current);

      // Expand window before showing item
      updateWindowSize(true);

      setItem(data);
      setIsScanning(false);
      setIsScanningFailed(false);
      setShowItem(true);

      fadeTimeoutRef.current = window.setTimeout(() => {
        setShowItem(false);
      }, 1500);

      clearTimeoutRef.current = window.setTimeout(() => {
        setItem(undefined);
        // Collapse window after item is cleared
        updateWindowSize(false);
      }, 2300);
    });

    EventsOn("scan-started", () => {
      setIsScanning(true);
      setIsScanningFailed(false);
    });

    EventsOn("scan-failed", () => {
      // Expand window to show failure message
      updateWindowSize(true);
      setIsScanningFailed(true);
      window.setTimeout(() => {
        setIsScanningFailed(false);
        // Collapse window after message disappears
        updateWindowSize(false);
      }, 2000);
    });

    EventsOn("toggle-visibility", () => {
      setIsVisible((prev) => !prev);
    });
  }, []);

  return (
    <div id="app">
      <div
        style={{
          display: "flex",
          alignItems: "end",
          flexDirection: "column",
          opacity: isVisible ? 1 : 0,
          transition: "opacity 0.3s ease-in-out",
        }}
      >
        {isScanning ? (
          <div>Scanning...</div>
        ) : isScanningFailed ? (
          <div>Scanning failed</div>
        ) : item ? (
          <>
            <div className="item-card fade-in">
              <img
                style={{
                  width: "100%",
                  height: "auto",
                  maxHeight: "70%",
                  objectFit: "contain",
                }}
                src={item.icon}
              />
              <span style={{ fontSize: "0.8rem" }}> $ {item.value}</span>
            </div>
            {(() => {
              const category = getItemCategory(item.id);
              const workshops = getWorkshopRequirements(item.id);
              const tags = getItemTags(item.id);

              if (category === "unknown" && workshops.length === 0 && tags.length === 0) return null;

              return (
                <>
                  {category !== "unknown" && (
                    <div className={`badge ${category === "keep" ? "unsafeToRecycle" : "safeToRecycle"}`}>
                      <span style={{ fontSize: "0.8rem" }}>
                        {category === "keep" ? "Keep" : "Recycle"}
                      </span>
                    </div>
                  )}
                  {tags.map((tag, index) => (
                    <div key={`tag-${index}`} className={`badge ${tag.toLowerCase()}`}>
                      <span style={{ fontSize: "0.65rem" }}>
                        {tag}
                      </span>
                    </div>
                  ))}
                  {workshops.map((req, index) => (
                    <div key={`ws-${index}`} className="badge workshop">
                      <span style={{ fontSize: "0.65rem" }}>
                        {req.workshop} L{req.level}
                      </span>
                    </div>
                  ))}
                </>
              );
            })()}
          </>
        ) : null}
      </div>
    </div>
  );
}

export default App;
