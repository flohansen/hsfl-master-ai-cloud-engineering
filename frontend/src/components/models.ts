
export interface ResponsePage<T> {
  page: Pagination;
  records: T[];
}

export interface Pagination {
  currentPage: number;
  pageSize: number;
  totalPages: number;
  totalRecords: number;
}

export interface BulletinBoardEntry {
  id: string;
  title: string;
  content: string;
  createdAt: string;
}
