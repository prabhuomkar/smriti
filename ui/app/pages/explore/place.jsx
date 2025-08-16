import { useParams } from "react-router";

export function meta() {
  return [{ title: "(TODO) Place Name - Smriti" }];
}

export default function Place() {
  const { id } = useParams();

  return <p>Place: {id}</p>;
}
