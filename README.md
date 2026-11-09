# dbgap-prep

A Pennsieve analysis workflow application to create dbGaP submission files.

## workflow

This application should be run in a workflow where it is preceded by a Data Source node to allow the user to choose the
dataset's `dataset_description.xlsx`, `subjects.xlsx`, and `samples.xlsx` as input.
This application should be followed in the workflow by a Data Target node with target type Pennsieve Dataset so that the
dbGaP submission files created by this application are uploaded back into the dataset being prepared for submission.

Documentation on the format and contents of the dbGaP submission files can be found
here: https://www.ncbi.nlm.nih.gov/gap/docs/submissionguide.

## app.yml

`app.yml` is the app's manifest: what the App Store shows and what a workflow needs to run it (its parameters, its
input and output, its compute) and the inputs a release is reviewed with before it is made public (`review`). The
schema is https://schemas.pennsieve.io/app-manifest/v1.json.

The parameters reach the app as environment variables of the same name, with the input files in `INPUT_DIR` and the
dbGaP files written to `OUTPUT_DIR`.

## runtimes

The image runs the same way as an ECS task and as a Lambda step: on Lambda, the platform's Lambda step runtime runs
the image's entrypoint with the same environment variables, so the app has no Lambda-specific code. Pennsieve decides
whether a release can run on Lambda from its `runtime` and its image; `computeTypes` only needs `standard`.

## publishing

Create a GitHub release (e.g. `v1.2.0`). Pennsieve builds the image and the version is in the App Store when the
build is done.

## running locally

If you have example input files you'd like to test with, run

```aiignore
 % go run cmd/local/main.go -i <input directory containing dataset_description.xlsx, subjects.xlsx, and samples.xlsx> -o <output directory> -c <consent group, GRU or HMB> -a <analyte type, DNA, RNA, or DNA/RNA>
```

`review-input/` holds a small example input (the one releases are reviewed with) that produces every dbGaP file:

```aiignore
 % go run cmd/local/main.go -i review-input -o out -c GRU -a DNA
```
