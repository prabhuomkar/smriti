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
    'about',
    'features',
    {
      type: 'category',
      label: 'Developer Guide',
      items: [
        'dev-guide/getting-started',
        'dev-guide/design-document',
        'dev-guide/folder-structure',
        'dev-guide/environment',
        'dev-guide/development',
        'dev-guide/testing',
        'dev-guide/roadmap',
        'dev-guide/notes',
      ],
    },
  ],
};

module.exports = sidebars;
