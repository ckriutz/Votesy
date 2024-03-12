import Image from "next/image";
import Header from "./Header";
import QuestionTable from "./QuestionTable";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center p-10">
      <Header />
      <QuestionTable />



    </main>
  );
}
