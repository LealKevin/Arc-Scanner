import { useEffect, useState, useRef } from "react";
import "./App.css";
import { EventsOn } from "../wailsjs/runtime/runtime";

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

  useEffect(() => {
    EventsOn("item-found", (data) => {
      if (fadeTimeoutRef.current) window.clearTimeout(fadeTimeoutRef.current);
      if (clearTimeoutRef.current) window.clearTimeout(clearTimeoutRef.current);

      setItem(data);
      setIsScanning(false);
      setIsScanningFailed(false);
      setShowItem(true);

      fadeTimeoutRef.current = window.setTimeout(() => {
        setShowItem(false);
      }, 1500);

      clearTimeoutRef.current = window.setTimeout(() => {
        setItem(undefined);
      }, 2300);
    });

    EventsOn("scan-started", () => {
      setIsScanning(true);
      setIsScanningFailed(false);
    });

    EventsOn("scan-failed", () => {
      setIsScanningFailed(true);
      window.setTimeout(() => {
        setIsScanningFailed(false);
      }, 2000);
    });

    EventsOn("toggle-visibility", () => {
      setIsVisible((prev) => !prev);
    });
  }, []);

  return (
    <div id="App">
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
            <div
              style={{
                width: "5rem",
                height: "5rem",
                display: "flex",
                alignItems: "center",
                flexDirection: "column",
                border: "1px solid white",
                justifyContent: "center",
                borderRadius: "0.5rem",
                backgroundColor: "rgba(0, 0, 0, 0.5)",
                overflow: "hidden",
                opacity: showItem ? 1 : 0,
                transition: "opacity 0.3s ease-in-out",
              }}
            >
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
          </>
        ) : null}
      </div>
    </div>
  );
}

export default App;
