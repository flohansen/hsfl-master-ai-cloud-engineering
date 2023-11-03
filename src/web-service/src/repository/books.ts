import axios from "axios";

export type Book = {
  Id: number;
  description: string;
};

export const getAllBooks = async () => {
  const response = await axios.get("/api/v1/books");
  return response.data as Book[];
};
