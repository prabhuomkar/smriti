import { useParams } from "react-router";

export function meta() {
  return [{ title: "(TODO) Shared Album Name - Smriti" }];
}

export default function SharedAlbum() {
  const { id } = useParams();

  return <p>Shared Album: {id}</p>;
}
