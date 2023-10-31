import { ModeToggle } from "@/components/mode-toggle.tsx";
import {useQuery} from "@tanstack/react-query";
import {Book, getAllBooks} from "@/repository/books.ts";

const BookCard = ({book}: {book: Book}) => {
    return <div>{book.description}</div>
}

function App() {
    const {data, isError, isLoading, isSuccess, error} = useQuery({
        queryKey: ["books"],
        queryFn: getAllBooks
    })

    if (isLoading) {
        return <div>Loading...</div>
    }

    if (isError) {
        return <div>Error {error.message}</div>
    }

    if (!isSuccess) {
        return <div>Something went wrong!</div>
    }

  return (
    <div>
      <ModeToggle />
        <div>{data.map(book => <BookCard book={book} />)}</div>
        <div>Hello World!</div>
    </div>
  );
}

export default App;
