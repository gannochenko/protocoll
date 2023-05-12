This 100% works:

~~~
name: Build and Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
      - name: Build
        run: make build
      - name: Create changelog
        id: github_release
        uses: mikepenz/release-changelog-builder-action@v3
        env:
          GITHUB_TOKEN: ${{ secrets.RELEASE_TOKEN }}
      - name: Create release
        uses: ncipollo/release-action@v1
        with:
          skipIfReleaseExists: true
          artifacts: "out/*"
          body: ${{steps.github_release.outputs.changelog}}
          token: ${{ secrets.RELEASE_TOKEN }}
      - uses: actions/checkout@v3
        with:
          ref: inst
          token: ${{ secrets.RELEASE_TOKEN }}
      - run: |
          date > generated.txt
          git config user.name github-actions
          git config user.email github-actions@github.com
          git add .
          git commit -m "generated2"
          git push origin inst
~~~