type Props = {
  isScanning: boolean;
  isFailed: boolean;
};

export function ScanStatus({ isScanning, isFailed }: Props) {
  if (isScanning) {
    return <div className="scan-status">Scanning...</div>;
  }
  if (isFailed) {
    return <div className="scan-status failed">Scanning failed</div>;
  }
  return null;
}
