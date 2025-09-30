# dbgap-prep
A Pennsieve analysis workflow application to create dbGaP submission files.

## workflow
This application should be run in a workflow where it is preceeded by https://github.com/Pennsieve/processor-pre-packages-v2 which will place the files `subjects.xlsx` and `samples.xlsx` into its input directory. This application should be followed by https://github.com/Pennsieve/processor-post-agent-v2 in the workflow so that the dbGaP submission files created by this application are uploaded back into the dataset being prepared for submission.

Documentation on the format and contents of the dbGaP submission files can be found here: https://www.ncbi.nlm.nih.gov/gap/docs/submissionguide.

