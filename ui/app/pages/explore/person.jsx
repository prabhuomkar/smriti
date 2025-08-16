import { useParams } from "react-router";

export function meta() {
  return [{ title: "(TODO) Person Name - Smriti" }];
}

export default function Person() {
  const { id } = useParams();

  return <p>Person: {id}</p>;
}
