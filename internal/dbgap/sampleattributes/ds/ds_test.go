package ds

import (
	"testing"

	"github.com/pennsieve/dbgap-prep/internal/dbgap/dd"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/ds"
	"github.com/pennsieve/dbgap-prep/internal/dbgap/sampleattributes/models"
	"github.com/pennsieve/dbgap-prep/internal/enums/analytetype"
	"github.com/pennsieve/dbgap-prep/internal/enums/istumor"
	"github.com/pennsieve/dbgap-prep/internal/samples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func TestWrite(t *testing.T) {
	writer := ds.NewXLSXWriter(t.TempDir(), DefaultFileNameBase)

	require.NoError(t, Write(writer, analytetype.DNA, istumor.YES, "https://doi.example.com/123/abc", models.CanonicalVariables, consentedSubjectSamples))

	actualFile, err := excelize.OpenFile(writer.Path())
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, actualFile.Close())
	}()

	records, err := actualFile.GetRows("Sheet1")
	require.NoError(t, err)

	// add one for the header
	assert.Len(t, records, expectedDataRowsCount+1)

	variableNames := dd.VariableNames(models.CanonicalVariables)
	assert.Equal(t, variableNames, records[0])
	assert.NotEqual(t, variableNames, records[1])

	analyteTypeIdx := indexOf(t, records[0], models.AnalyteTypeVar.Name)
	isTumorIdx := indexOf(t, records[0], models.IsTumorVar.Name)
	sparcDOIURLIdx := indexOf(t, records[0], models.SPARCDatasetDOIVar.Name)
	for _, dataRow := range records[1:] {
		assert.Equal(t, analytetype.DNA.String(), dataRow[analyteTypeIdx])
		assert.Equal(t, istumor.YES.String(), dataRow[isTumorIdx])
		assert.Equal(t, "https://doi.example.com/123/abc", dataRow[sparcDOIURLIdx])
	}
}

func indexOf(t *testing.T, header []string, name string) int {
	t.Helper()
	for i, h := range header {
		if h == name {
			return i
		}
	}
	require.Failf(t, "variable not found in header", "%q not found in %v", name, header)
	return -1
}

func TestNewToRow(t *testing.T) {
	rowVariables := []dd.Variable{*dd.SampleIDVar, models.AnalyteTypeVar, models.IsTumorVar, models.SPARCDatasetDOIVar}
	sampleIDIdx := indexOf(t, dd.VariableNames(rowVariables), dd.SampleIDVar.Name)
	analyteTypeIdx := indexOf(t, dd.VariableNames(rowVariables), models.AnalyteTypeVar.Name)
	isTumorIdx := indexOf(t, dd.VariableNames(rowVariables), models.IsTumorVar.Name)
	sparcDOIURLIdx := indexOf(t, dd.VariableNames(rowVariables), models.SPARCDatasetDOIVar.Name)

	sample := samples.Sample{ID: "sam-1", SubjectID: "sub-1", Values: map[string]string{}}

	for _, analyteType := range []analytetype.Type{analytetype.DNA, analytetype.RNA, analytetype.DNARNA} {
		for _, isTumor := range []istumor.Value{istumor.YES, istumor.NO} {
			t.Run(analyteType.String()+"/"+isTumor.String(), func(t *testing.T) {
				sparcDOIURL := "https://doi.example.com/123/abc"
				row := NewToRow(analyteType, isTumor, sparcDOIURL)(rowVariables, sample)

				assert.Equal(t, sample.ID, row[sampleIDIdx])
				assert.Equal(t, analyteType.String(), row[analyteTypeIdx])
				assert.Equal(t, isTumor.String(), row[isTumorIdx])
				assert.Equal(t, sparcDOIURL, row[sparcDOIURLIdx])
			})
		}
	}
}

var consentedSubjectSamples = []samples.Sample{
	samples.Sample{ID: "sam-hSG-lib-4", SubjectID: "sub-hSG-Female-1", Values: map[string]string{"RNA concentration": "", "RNA concentration method": "", "RNA purity": "", "RNA quality": "", "RNA quality method": "", "also in dataset": "", "amputation": "not-asked", "body temp": "", "cell type": "", "collection institution": "WashU", "cross clamp (first)": "", "cross clamp (last)": "", "cross clamp time": "", "date of derivation": "", "donor status": "", "experimental log file path": "", "fixation method": "Fresh Frozen", "fixation temp": "-170", "fixation time": "", "freeze thaw cycles": "", "freezing method": "", "freezing temp": "", "incision time": "", "ischemic time": "", "laboratory internal id": "hSG-4", "laterality": "", "member of": "", "metadata only": "TRUE", "nCells": "6882", "nCount_ATAC": "4880.5233943621", "nCount_RNA": "3662.77651845394", "nFeature_ATAC": "2382.46875908166", "nFeature_RNA": "1776.60665504214", "number of directly derived samples": "", "pathology": "", "plane of section": "", "pool id": "", "post-mortem interval": "", "protein concentration": "", "protein concentration method": "", "protocol title": "", "protocol url or doi": "", "reference atlas": "", "sample anatomical location": "sympathetic ganglion", "sample collection site": "", "sample experimental group": "", "sample type": "tissue", "storage temp": "-80", "time of sample collection": "", "was derived from": ""}},
	samples.Sample{ID: "sam-hSG-lib-7", SubjectID: "sub-hSG-Female-2", Values: map[string]string{"RNA concentration": "", "RNA concentration method": "", "RNA purity": "", "RNA quality": "", "RNA quality method": "", "also in dataset": "", "amputation": "not-asked", "body temp": "", "cell type": "", "collection institution": "WashU", "cross clamp (first)": "", "cross clamp (last)": "", "cross clamp time": "", "date of derivation": "", "donor status": "", "experimental log file path": "", "fixation method": "Fresh Frozen", "fixation temp": "-170", "fixation time": "", "freeze thaw cycles": "", "freezing method": "", "freezing temp": "", "incision time": "", "ischemic time": "", "laboratory internal id": "hSG-7", "laterality": "", "member of": "", "metadata only": "TRUE", "nCells": "5951", "nCount_ATAC": "12078.1188035624", "nCount_RNA": "3806.49000168039", "nFeature_ATAC": "5239.86304822719", "nFeature_RNA": "1739.59519408503", "number of directly derived samples": "", "pathology": "", "plane of section": "", "pool id": "", "post-mortem interval": "", "protein concentration": "", "protein concentration method": "", "protocol title": "", "protocol url or doi": "", "reference atlas": "", "sample anatomical location": "sympathetic ganglion", "sample collection site": "", "sample experimental group": "", "sample type": "tissue", "storage temp": "-80", "time of sample collection": "", "was derived from": ""}}, samples.Sample{ID: "sam-hSG-lib-8", SubjectID: "sub-hSG-Female-2", Values: map[string]string{"RNA concentration": "", "RNA concentration method": "", "RNA purity": "", "RNA quality": "", "RNA quality method": "", "also in dataset": "", "amputation": "not-asked", "body temp": "", "cell type": "", "collection institution": "WashU", "cross clamp (first)": "", "cross clamp (last)": "", "cross clamp time": "", "date of derivation": "", "donor status": "", "experimental log file path": "", "fixation method": "Fresh Frozen", "fixation temp": "-170", "fixation time": "", "freeze thaw cycles": "", "freezing method": "", "freezing temp": "", "incision time": "", "ischemic time": "", "laboratory internal id": "hSG-8", "laterality": "", "member of": "", "metadata only": "TRUE", "nCells": "5678", "nCount_ATAC": "12090.1377245509", "nCount_RNA": "3724.34114124692", "nFeature_ATAC": "5230.64864388869", "nFeature_RNA": "1712.71556886227", "number of directly derived samples": "", "pathology": "", "plane of section": "", "pool id": "", "post-mortem interval": "", "protein concentration": "", "protein concentration method": "", "protocol title": "", "protocol url or doi": "", "reference atlas": "", "sample anatomical location": "sympathetic ganglion", "sample collection site": "", "sample experimental group": "", "sample type": "tissue", "storage temp": "-80", "time of sample collection": "", "was derived from": ""}},
	samples.Sample{ID: "sam-D-20210810A-SG4-CURLS-02", SubjectID: "sub-hSG-Male-1", Values: map[string]string{"RNA concentration": "", "RNA concentration method": "", "RNA purity": "", "RNA quality": "", "RNA quality method": "", "also in dataset": "", "amputation": "", "body temp": "", "cell type": "", "collection institution": "", "cross clamp (first)": "", "cross clamp (last)": "", "cross clamp time": "", "date of derivation": "", "donor status": "", "experimental log file path": "", "fixation method": "", "fixation temp": "", "fixation time": "", "freeze thaw cycles": "", "freezing method": "", "freezing temp": "", "incision time": "", "ischemic time": "", "laboratory internal id": "", "laterality": "", "member of": "", "metadata only": "TRUE", "nCells": "", "nCount_ATAC": "", "nCount_RNA": "", "nFeature_ATAC": "", "nFeature_RNA": "", "number of directly derived samples": "", "pathology": "", "plane of section": "", "pool id": "", "post-mortem interval": "", "protein concentration": "", "protein concentration method": "", "protocol title": "", "protocol url or doi": "", "reference atlas": "", "sample anatomical location": "sympathetic ganglion", "sample collection site": "", "sample experimental group": "", "sample type": "tissue", "storage temp": "", "time of sample collection": "", "was derived from": ""}}, samples.Sample{ID: "sam-hSG-lib-2", SubjectID: "sub-hSG-Male-1", Values: map[string]string{"RNA concentration": "", "RNA concentration method": "", "RNA purity": "", "RNA quality": "", "RNA quality method": "", "also in dataset": "", "amputation": "not-asked", "body temp": "", "cell type": "", "collection institution": "WashU", "cross clamp (first)": "", "cross clamp (last)": "", "cross clamp time": "", "date of derivation": "", "donor status": "", "experimental log file path": "", "fixation method": "Fresh Frozen", "fixation temp": "-170", "fixation time": "", "freeze thaw cycles": "", "freezing method": "", "freezing temp": "", "incision time": "", "ischemic time": "", "laboratory internal id": "hSG-2", "laterality": "", "member of": "", "metadata only": "TRUE", "nCells": "8189", "nCount_ATAC": "6925.28257418488", "nCount_RNA": "4138.48540725363", "nFeature_ATAC": "3266.53730614239", "nFeature_RNA": "1986.50799853462", "number of directly derived samples": "", "pathology": "", "plane of section": "", "pool id": "", "post-mortem interval": "", "protein concentration": "", "protein concentration method": "", "protocol title": "", "protocol url or doi": "", "reference atlas": "", "sample anatomical location": "sympathetic ganglion", "sample collection site": "", "sample experimental group": "", "sample type": "tissue", "storage temp": "-80", "time of sample collection": "", "was derived from": ""}},
	samples.Sample{ID: "sam-hSG-lib-3", SubjectID: "sub-hSG-Male-2", Values: map[string]string{"RNA concentration": "", "RNA concentration method": "", "RNA purity": "", "RNA quality": "", "RNA quality method": "", "also in dataset": "", "amputation": "not-asked", "body temp": "", "cell type": "", "collection institution": "WashU", "cross clamp (first)": "", "cross clamp (last)": "", "cross clamp time": "", "date of derivation": "", "donor status": "", "experimental log file path": "", "fixation method": "Fresh Frozen", "fixation temp": "-170", "fixation time": "", "freeze thaw cycles": "", "freezing method": "", "freezing temp": "", "incision time": "", "ischemic time": "", "laboratory internal id": "hSG-3", "laterality": "", "member of": "", "metadata only": "TRUE", "nCells": "6465", "nCount_ATAC": "3825.2909512761", "nCount_RNA": "3426.73302397525", "nFeature_ATAC": "1941.78128383604", "nFeature_RNA": "1741.47981438515", "number of directly derived samples": "", "pathology": "", "plane of section": "", "pool id": "", "post-mortem interval": "", "protein concentration": "", "protein concentration method": "", "protocol title": "", "protocol url or doi": "", "reference atlas": "", "sample anatomical location": "sympathetic ganglion", "sample collection site": "", "sample experimental group": "", "sample type": "tissue", "storage temp": "-80", "time of sample collection": "", "was derived from": ""}},
	samples.Sample{ID: "sam-hSG-lib-5", SubjectID: "sub-hSG-Male-3", Values: map[string]string{"RNA concentration": "", "RNA concentration method": "", "RNA purity": "", "RNA quality": "", "RNA quality method": "", "also in dataset": "", "amputation": "not-asked", "body temp": "", "cell type": "", "collection institution": "WashU", "cross clamp (first)": "", "cross clamp (last)": "", "cross clamp time": "", "date of derivation": "", "donor status": "", "experimental log file path": "", "fixation method": "Fresh Frozen", "fixation temp": "-170", "fixation time": "", "freeze thaw cycles": "", "freezing method": "", "freezing temp": "", "incision time": "", "ischemic time": "", "laboratory internal id": "hSG-5", "laterality": "", "member of": "", "metadata only": "TRUE", "nCells": "4905", "nCount_ATAC": "15876.6799184506", "nCount_RNA": "4260.05626911315", "nFeature_ATAC": "6432.90948012232", "nFeature_RNA": "1823.30132517839", "number of directly derived samples": "", "pathology": "", "plane of section": "", "pool id": "", "post-mortem interval": "", "protein concentration": "", "protein concentration method": "", "protocol title": "", "protocol url or doi": "", "reference atlas": "", "sample anatomical location": "sympathetic ganglion", "sample collection site": "", "sample experimental group": "", "sample type": "tissue", "storage temp": "-80", "time of sample collection": "", "was derived from": ""}}, samples.Sample{ID: "sam-hSG-lib-6", SubjectID: "sub-hSG-Male-3", Values: map[string]string{"RNA concentration": "", "RNA concentration method": "", "RNA purity": "", "RNA quality": "", "RNA quality method": "", "also in dataset": "", "amputation": "not-asked", "body temp": "", "cell type": "", "collection institution": "WashU", "cross clamp (first)": "", "cross clamp (last)": "", "cross clamp time": "", "date of derivation": "", "donor status": "", "experimental log file path": "", "fixation method": "Fresh Frozen", "fixation temp": "-170", "fixation time": "", "freeze thaw cycles": "", "freezing method": "", "freezing temp": "", "incision time": "", "ischemic time": "", "laboratory internal id": "hSG-6", "laterality": "", "member of": "", "metadata only": "TRUE", "nCells": "4764", "nCount_ATAC": "12436.6511335013", "nCount_RNA": "4624.95046179681", "nFeature_ATAC": "5285.26994122586", "nFeature_RNA": "1936.36796809404", "number of directly derived samples": "", "pathology": "", "plane of section": "", "pool id": "", "post-mortem interval": "", "protein concentration": "", "protein concentration method": "", "protocol title": "", "protocol url or doi": "", "reference atlas": "", "sample anatomical location": "sympathetic ganglion", "sample collection site": "", "sample experimental group": "", "sample type": "tissue", "storage temp": "-80", "time of sample collection": "", "was derived from": ""}},
	samples.Sample{ID: "sam-D-20210629A-SG-CURLS-02", SubjectID: "sub-hSG-Male-4", Values: map[string]string{"RNA concentration": "", "RNA concentration method": "", "RNA purity": "", "RNA quality": "", "RNA quality method": "", "also in dataset": "", "amputation": "", "body temp": "", "cell type": "", "collection institution": "", "cross clamp (first)": "", "cross clamp (last)": "", "cross clamp time": "", "date of derivation": "", "donor status": "", "experimental log file path": "", "fixation method": "", "fixation temp": "", "fixation time": "", "freeze thaw cycles": "", "freezing method": "", "freezing temp": "", "incision time": "", "ischemic time": "", "laboratory internal id": "", "laterality": "", "member of": "", "metadata only": "TRUE", "nCells": "", "nCount_ATAC": "", "nCount_RNA": "", "nFeature_ATAC": "", "nFeature_RNA": "", "number of directly derived samples": "", "pathology": "", "plane of section": "", "pool id": "", "post-mortem interval": "", "protein concentration": "", "protein concentration method": "", "protocol title": "", "protocol url or doi": "", "reference atlas": "", "sample anatomical location": "sympathetic ganglion", "sample collection site": "", "sample experimental group": "", "sample type": "tissue", "storage temp": "", "time of sample collection": "", "was derived from": ""}},
	samples.Sample{ID: "sam-D-20220418A-SG-10A-CURLS-02", SubjectID: "sub-hSG-Male-5", Values: map[string]string{"RNA concentration": "", "RNA concentration method": "", "RNA purity": "", "RNA quality": "", "RNA quality method": "", "also in dataset": "", "amputation": "", "body temp": "", "cell type": "", "collection institution": "", "cross clamp (first)": "", "cross clamp (last)": "", "cross clamp time": "", "date of derivation": "", "donor status": "", "experimental log file path": "", "fixation method": "", "fixation temp": "", "fixation time": "", "freeze thaw cycles": "", "freezing method": "", "freezing temp": "", "incision time": "", "ischemic time": "", "laboratory internal id": "", "laterality": "", "member of": "", "metadata only": "TRUE", "nCells": "", "nCount_ATAC": "", "nCount_RNA": "", "nFeature_ATAC": "", "nFeature_RNA": "", "number of directly derived samples": "", "pathology": "", "plane of section": "", "pool id": "", "post-mortem interval": "", "protein concentration": "", "protein concentration method": "", "protocol title": "", "protocol url or doi": "", "reference atlas": "", "sample anatomical location": "sympathetic ganglion", "sample collection site": "", "sample experimental group": "", "sample type": "tissue", "storage temp": "", "time of sample collection": "", "was derived from": ""}}}

// expectedDataRowsCount is the expected number of non-header rows. One for each sample in consentedSubjectSamples
var expectedDataRowsCount = len(consentedSubjectSamples)
