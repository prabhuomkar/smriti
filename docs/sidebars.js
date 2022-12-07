/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  devGuideSidebar: [
    {
      type: "category",
      label: "User Guide",
      items: [
        "user-guide/introduction",
        "user-guide/installation",
        {
          type: "category",
          label: "Features",
          items: [
            "user-guide/features/photos",
            "user-guide/features/albums",
            "user-guide/features/users",
            "user-guide/features/metadata",
            "user-guide/features/places",
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
