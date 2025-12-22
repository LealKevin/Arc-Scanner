import { useState, useEffect } from "react";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { DownloadUpdate, ApplyUpdateAndRestart } from "../../wailsjs/go/main/App";
import type { UpdateInfo } from "../types";
import "./UpdateNotification.css";

type UpdateState = "available" | "downloading" | "ready" | "error";

export function UpdateNotification() {
  const [updateInfo, setUpdateInfo] = useState<UpdateInfo | null>(null);
  const [state, setState] = useState<UpdateState>("available");
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const unsubAvailable = EventsOn("update-available", (info: UpdateInfo) => {
      setUpdateInfo(info);
      setState("available");
    });

    const unsubProgress = EventsOn("update-progress", (percent: number) => {
      setProgress(percent);
    });

    const unsubReady = EventsOn("update-ready", () => {
      setState("ready");
    });

    const unsubError = EventsOn("update-error", (errMsg: string) => {
      setState("error");
      setError(errMsg);
    });

    return () => {
      unsubAvailable();
      unsubProgress();
      unsubReady();
      unsubError();
    };
  }, []);

  const handleDownload = async () => {
    if (!updateInfo) return;
    setState("downloading");
    setProgress(0);
    try {
      await DownloadUpdate(updateInfo);
    } catch (err) {
      setState("error");
      setError(String(err));
    }
  };

  const handleRestart = async () => {
    try {
      await ApplyUpdateAndRestart();
    } catch (err) {
      setState("error");
      setError(String(err));
    }
  };

  const handleDismiss = () => {
    setUpdateInfo(null);
  };

  if (!updateInfo) return null;

  return (
    <div className="update-notification">
      {state === "available" && (
        <>
          <span className="update-text">
            v{updateInfo.version} available
          </span>
          <button className="update-btn" onClick={handleDownload}>
            Update
          </button>
          <button className="dismiss-btn" onClick={handleDismiss}>
            x
          </button>
        </>
      )}

      {state === "downloading" && (
        <>
          <span className="update-text">Downloading...</span>
          <div className="progress-bar">
            <div className="progress-fill" style={{ width: `${progress}%` }} />
          </div>
          <span className="progress-text">{progress}%</span>
        </>
      )}

      {state === "ready" && (
        <>
          <span className="update-text">Ready to install</span>
          <button className="update-btn" onClick={handleRestart}>
            Restart
          </button>
        </>
      )}

      {state === "error" && (
        <>
          <span className="update-text error">Update failed</span>
          <button className="update-btn" onClick={handleDownload}>
            Retry
          </button>
          <button className="dismiss-btn" onClick={handleDismiss}>
            x
          </button>
        </>
      )}
    </div>
  );
}
