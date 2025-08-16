import { useParams } from "react-router";

export function meta() {
  return [{ title: "(TODO) Album Name - Smriti" }];
}

export default function Album() {
  const { id } = useParams();

  return <p>Album: {id}</p>;
}
