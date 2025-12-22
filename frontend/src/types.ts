export type Item = {
  id: string;
  name: string;
  value: number;
  icon: string;
};

export type ItemFoundEvent = Item;

export type UpdateInfo = {
  version: string;
  url: string;
  downloadUrl: string;
  releaseNotes: string;
  publishedAt: string;
};
