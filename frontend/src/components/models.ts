
export interface ResponsePage<T> {
  totalPages: number;
  currentPage: number;
  pageSize: number;
  records: T[];
}
export interface BulletinBoardEntry {
  id: string;
  title: string;
  description: string;
  createdAt: string;
}
