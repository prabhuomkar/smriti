import { useParams } from "react-router";

export function meta() {
  return [{ title: "(TODO) Thing Name - Smriti" }];
}

export default function Thing() {
  const { id } = useParams();

  return <p>Thing: {id}</p>;
}
