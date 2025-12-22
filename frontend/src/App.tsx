import { useEffect, useState, useCallback } from "react";
import "./App.css";
import {
  EventsOn,
  WindowSetSize,
  WindowSetPosition,
  ScreenGetAll,
} from "../wailsjs/runtime/runtime";
import type { Item, ItemFoundEvent } from "./types";
import { useTimeout } from "./hooks/useTimeout";
import { ScanStatus } from "./components/ScanStatus";
import { ItemCard } from "./components/ItemCard";
import { ItemBadges } from "./components/ItemBadges";

const WINDOW_WIDTH_VISIBLE = 100;
const WINDOW_HEIGHT_VISIBLE = 220;
const WINDOW_WIDTH_HIDDEN = 1;
const WINDOW_HEIGHT_HIDDEN = 1;

function App() {
  const [item, setItem] = useState<Item>();
  const [isScanning, setIsScanning] = useState(false);
  const [isScanningFailed, setIsScanningFailed] = useState(false);
  const [isVisible, setIsVisible] = useState(true);
  const [showItem, setShowItem] = useState(false);

  const fadeTimeout = useTimeout();
  const clearTimeout = useTimeout();
  const failedTimeout = useTimeout();

  const updateWindowSize = useCallback(async (visible: boolean) => {
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
  }, []);

  useEffect(() => {
    updateWindowSize(false);

    const handleItemFound = (data: ItemFoundEvent) => {
      fadeTimeout.clear();
      clearTimeout.clear();
      failedTimeout.clear();

      updateWindowSize(true);

      setItem(data);
      setIsScanning(false);
      setIsScanningFailed(false);
      setShowItem(true);

      fadeTimeout.set(() => {
        setShowItem(false);
        setItem(undefined);
        updateWindowSize(false);
      }, 1500);
    };

    const handleScanStarted = () => {
      fadeTimeout.clear();
      clearTimeout.clear();
      failedTimeout.clear();
      setItem(undefined);
      setShowItem(false);
      setIsScanning(true);
      setIsScanningFailed(false);
    };

    const handleScanFailed = () => {
      fadeTimeout.clear();
      clearTimeout.clear();
      setItem(undefined);
      setShowItem(false);
      updateWindowSize(true);
      setIsScanningFailed(true);
      failedTimeout.set(() => {
        setIsScanningFailed(false);
        updateWindowSize(false);
      }, 2000);
    };

    const handleToggleVisibility = () => {
      setIsVisible((prev) => !prev);
    };

    const unsubItemFound = EventsOn("item-found", handleItemFound);
    const unsubScanStarted = EventsOn("scan-started", handleScanStarted);
    const unsubScanFailed = EventsOn("scan-failed", handleScanFailed);
    const unsubToggle = EventsOn("toggle-visibility", handleToggleVisibility);

    return () => {
      unsubItemFound();
      unsubScanStarted();
      unsubScanFailed();
      unsubToggle();
    };
  }, [updateWindowSize, fadeTimeout, clearTimeout, failedTimeout]);

  return (
    <div id="app">
      <div className={`container ${isVisible ? "visible" : "hidden"}`}>
        <ScanStatus isScanning={isScanning} isFailed={isScanningFailed} />
        {item && (
          <>
            <ItemCard item={item} className={showItem ? "fade-in" : ""} />
            <ItemBadges itemId={item.id} className={showItem ? "fade-in" : ""} />
          </>
        )}
      </div>
    </div>
  );
}

export default App;
