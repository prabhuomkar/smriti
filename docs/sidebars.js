const sidebars = {
  devGuideSidebar: [
    {
      type: "category",
      label: "User Guide",
      items: [
        "user-guide/introduction",
        "user-guide/installation",
        "user-guide/deployment",
        {
          type: "category",
          label: "Features",
          items: [
            "user-guide/features/users",
            "user-guide/features/metadata",
            "user-guide/features/library",
            "user-guide/features/albums",
            "user-guide/features/places",
            "user-guide/features/things",
            "user-guide/features/people",
            "user-guide/features/explore",
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Developer Guide",
      items: [
        "dev-guide/introduction",
        "dev-guide/environment",
        "dev-guide/folder-structure",
        "dev-guide/contribution",
        "dev-guide/roadmap",
      ],
    },
  ],
};

module.exports = sidebars;
