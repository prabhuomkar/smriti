import { index, route, prefix } from "@react-router/dev/routes";

export default [
  ...prefix("admin", [
    index("pages/admin/index.jsx"),
    route("users", "pages/admin/users.jsx"),
  ]),
  index("pages/home/index.jsx"),
  ...prefix("auth", [route("login", "pages/auth/login.jsx")]),
  route("albums", "pages/albums/index.jsx"),
  route("album/:id", "pages/albums/album.jsx"),
  route("sharing", "pages/sharing/index.jsx"),
  route("share/:id", "pages/sharing/sharedalbum.jsx"),
  route("favourites", "pages/library/favourites.jsx"),
  route("hidden", "pages/library/hidden.jsx"),
  route("trash", "pages/library/trash.jsx"),
  route("things", "pages/explore/things.jsx"),
  route("thing/:id", "pages/explore/thing.jsx"),
  route("places", "pages/explore/places.jsx"),
  route("place/:id", "pages/explore/place.jsx"),
  route("people", "pages/explore/people.jsx"),
  route("person/:id", "pages/explore/person.jsx"),
];
