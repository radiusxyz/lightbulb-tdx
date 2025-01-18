package attest

import (
	tdxpb "github.com/google/go-tdx-guest/proto/tdx"
	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
)

// ConvertQuoteV4ToQuote converts a QuoteV4 object to a Quote object.
func ConvertQuoteV4ToQuote(qv4 *tdxpb.QuoteV4) *attestpb.Quote {
	if qv4 == nil {
		return nil
	}

	return &attestpb.Quote{
		Header:        convertHeader(qv4.Header),
		TdQuoteBody:   convertTDQuoteBody(qv4.TdQuoteBody),
		SignedDataSize: qv4.SignedDataSize,
		SignedData:    convertSignedData(qv4.SignedData),
		ExtraBytes:    qv4.ExtraBytes,
	}
}

func convertHeader(h *tdxpb.Header) *attestpb.Header {
	if h == nil {
		return nil
	}
	return &attestpb.Header{
		Version:            h.Version,
		AttestationKeyType: h.AttestationKeyType,
		TeeType:            h.TeeType,
		QeSvn:              h.QeSvn,
		PceSvn:             h.PceSvn,
		QeVendorId:         h.QeVendorId,
		UserData:           h.UserData,
	}
}

func convertTDQuoteBody(tb *tdxpb.TDQuoteBody) *attestpb.TDQuoteBody {
	if tb == nil {
		return nil
	}
	return &attestpb.TDQuoteBody{
		TeeTcbSvn:      tb.TeeTcbSvn,
		MrSeam:         tb.MrSeam,
		MrSignerSeam:   tb.MrSignerSeam,
		SeamAttributes: tb.SeamAttributes,
		TdAttributes:   tb.TdAttributes,
		Xfam:           tb.Xfam,
		MrTd:           tb.MrTd,
		MrConfigId:     tb.MrConfigId,
		MrOwner:        tb.MrOwner,
		MrOwnerConfig:  tb.MrOwnerConfig,
		Rtmrs:          tb.Rtmrs,
		ReportData:     tb.ReportData,
	}
}

func convertSignedData(sd *tdxpb.Ecdsa256BitQuoteV4AuthData) *attestpb.Ecdsa256BitQuoteV4AuthData {
	if sd == nil {
		return nil
	}
	return &attestpb.Ecdsa256BitQuoteV4AuthData{
		Signature:           sd.Signature,
		EcdsaAttestationKey: sd.EcdsaAttestationKey,
		CertificationData:   convertCertificationData(sd.CertificationData),
	}
}

func convertCertificationData(cd *tdxpb.CertificationData) *attestpb.CertificationData {
	if cd == nil {
		return nil
	}
	return &attestpb.CertificationData{
		CertificateDataType: cd.CertificateDataType,
		Size:                cd.Size,
		QeReportCertificationData: convertQEReportCertificationData(cd.QeReportCertificationData),
	}
}

func convertQEReportCertificationData(qe *tdxpb.QEReportCertificationData) *attestpb.QEReportCertificationData {
	if qe == nil {
		return nil
	}
	return &attestpb.QEReportCertificationData{
		QeReport:                convertEnclaveReport(qe.QeReport),
		QeReportSignature:       qe.QeReportSignature,
		QeAuthData:              convertQeAuthData(qe.QeAuthData),
		PckCertificateChainData: convertPckCertificateChainData(qe.PckCertificateChainData),
	}
}

func convertEnclaveReport(er *tdxpb.EnclaveReport) *attestpb.EnclaveReport {
	if er == nil {
		return nil
	}
	return &attestpb.EnclaveReport{
		CpuSvn:     er.CpuSvn,
		MiscSelect: er.MiscSelect,
		Reserved1:  er.Reserved1,
		Attributes: er.Attributes,
		MrEnclave:  er.MrEnclave,
		Reserved2:  er.Reserved2,
		MrSigner:   er.MrSigner,
		Reserved3:  er.Reserved3,
		IsvProdId:  er.IsvProdId,
		IsvSvn:     er.IsvSvn,
		Reserved4:  er.Reserved4,
		ReportData: er.ReportData,
	}
}

func convertQeAuthData(qa *tdxpb.QeAuthData) *attestpb.QeAuthData {
	if qa == nil {
		return nil
	}
	return &attestpb.QeAuthData{
		ParsedDataSize: qa.ParsedDataSize,
		Data:           qa.Data,
	}
}

func convertPckCertificateChainData(pcc *tdxpb.PCKCertificateChainData) *attestpb.PCKCertificateChainData {
	if pcc == nil {
		return nil
	}
	return &attestpb.PCKCertificateChainData{
		CertificateDataType: pcc.CertificateDataType,
		Size:                pcc.Size,
		PckCertChain:        pcc.PckCertChain,
	}
}
