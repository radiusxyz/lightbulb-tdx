package auction

import (
	"time"

	auctionpb "github.com/radiusxyz/lightbulb-tdx/proto/auction"
)

// Tx represents a transaction submitted by a bidder.
type Tx struct {
	TxData string // Raw transaction data.
}

// Bid represents a bid submitted by a buyer.
type Bid struct {
	ChainID         int64  // Unique identifier of the chain.
	AuctionID       string // Unique identifier of the auction.
	BidderAddr      string // Address of the bidder.
	BidAmount       int64  // Amount of the bid.
	BidderSignature string // Signature of the bidder.
	TxList          []Tx   // List of transactions associated with the bid.
}

// AuctionInfo represents the details of an auction.
type AuctionInfo struct {
	AuctionID       string    // Unique identifier of the auction.
	ChainID         int64     // Unique identifier of the chain.
	StartTime       time.Time // Start time of the auction.
	EndTime         time.Time // End time of the auction.
	SellerAddress   string    // Address of the seller.
	BlockNumber     int64     // Block number where the auction is registered.
	BlockspaceSize  int64     // Block space size being auctioned.
	SellerSignature string    // Seller's signature for the auction.
}

// AuctionState represents the current state of an auction.
type AuctionState struct {
	AuctionInfo  AuctionInfo // Details of the auction.
	BidList      []Bid       // List of all bids submitted.
	SortedTxList []Tx        // List of all transactions sorted.
	IsEnded      bool        // Indicates whether the auction has ended.
}

// ConvertProtobufBidToDomainBid converts a protobuf Bid to a domain Bid.
func ConvertProtobufBidToDomain(pbBid *auctionpb.Bid) Bid {
	var txList []Tx
	for _, pbTx := range pbBid.TxList {
		txList = append(txList, Tx{
			TxData: pbTx.GetTxData(),
		})
	}

	return Bid{
		ChainID:         pbBid.GetChainId(),
		AuctionID:       pbBid.GetAuctionId(),
		BidderAddr:      pbBid.GetBidderAddr(),
		BidAmount:       pbBid.GetBidAmount(),
		BidderSignature: pbBid.GetBidderSignature(),
		TxList:          txList,
	}
}

// ConvertProtobufBidsToDomainBids converts a slice of protobuf Bids to a slice of domain Bids.
func ConvertProtobufBidsToDomain(pbBids []*auctionpb.Bid) []Bid {
	var domainBids []Bid
	for _, pbBid := range pbBids {
		domainBids = append(domainBids, ConvertProtobufBidToDomain(pbBid))
	}
	return domainBids
}

func ConvertProtobufAuctionInfoToDomain(pbAuctionInfo *auctionpb.AuctionInfo) AuctionInfo {
	return AuctionInfo{
		AuctionID:       pbAuctionInfo.GetAuctionId(),
		ChainID:         pbAuctionInfo.GetChainId(),
		StartTime:       time.Unix(0, pbAuctionInfo.GetStartTime()*int64(time.Millisecond)),
		EndTime:         time.Unix(0, pbAuctionInfo.GetEndTime()*int64(time.Millisecond)),
		SellerAddress:   pbAuctionInfo.GetSellerAddress(),
		BlockNumber:     pbAuctionInfo.GetBlockNumber(),
		BlockspaceSize:  pbAuctionInfo.GetBlockspaceSize(),
		SellerSignature: pbAuctionInfo.GetSellerSignature(),
	}
}

func ConvertDomainAuctionInfoToProtobuf(domainAuctionInfo AuctionInfo) *auctionpb.AuctionInfo {
	return &auctionpb.AuctionInfo{
		AuctionId:       domainAuctionInfo.AuctionID,
		ChainId:         domainAuctionInfo.ChainID,
		StartTime:       domainAuctionInfo.StartTime.UnixNano() / int64(time.Millisecond),
		EndTime:         domainAuctionInfo.EndTime.UnixNano() / int64(time.Millisecond),
		SellerAddress:   domainAuctionInfo.SellerAddress,
		BlockNumber:     domainAuctionInfo.BlockNumber,
		BlockspaceSize:  domainAuctionInfo.BlockspaceSize,
		SellerSignature: domainAuctionInfo.SellerSignature,
	}
}

func ConvertDomainBidToProtobuf(domainBid Bid) *auctionpb.Bid {
	var pbTxList []*auctionpb.Tx
	for _, tx := range domainBid.TxList {
		pbTxList = append(pbTxList, &auctionpb.Tx{
			TxData: tx.TxData,
		})
	}

	return &auctionpb.Bid{
		ChainId:         domainBid.ChainID,
		AuctionId:       domainBid.AuctionID,
		BidderAddr:      domainBid.BidderAddr,
		BidAmount:       domainBid.BidAmount,
		BidderSignature: domainBid.BidderSignature,
		TxList:          pbTxList,
	}
}

func ConvertDomainBidsToProtobuf(domainBids []Bid) []*auctionpb.Bid {
	var pbBids []*auctionpb.Bid
	for _, domainBid := range domainBids {
		pbBids = append(pbBids, ConvertDomainBidToProtobuf(domainBid))
	}
	return pbBids
}

func ConvertDomainTxToProtobuf(domainTx Tx) *auctionpb.Tx {
	return &auctionpb.Tx{
		TxData: domainTx.TxData,
	}
}

func ConvertDomainTxsToProtobuf(domainTxs []Tx) []*auctionpb.Tx {
	var pbTxs []*auctionpb.Tx
	for _, domainTx := range domainTxs {
		pbTxs = append(pbTxs, ConvertDomainTxToProtobuf(domainTx))
	}
	return pbTxs
}

func ConvertProtobufTxToDomain(pbTx *auctionpb.Tx) Tx {
	return Tx{
		TxData: pbTx.GetTxData(),
	}
}

func ConvertProtobufTxsToDomain(pbTxs []*auctionpb.Tx) []Tx {
	var domainTxs []Tx
	for _, pbTx := range pbTxs {
		domainTxs = append(domainTxs, ConvertProtobufTxToDomain(pbTx))
	}
	return domainTxs
}

func ConvertProtobufAuctionStateToDomain(pbAuctionState *auctionpb.AuctionState) AuctionState {
	return AuctionState{
		AuctionInfo:  ConvertProtobufAuctionInfoToDomain(pbAuctionState.GetAuctionInfo()),
		BidList:      ConvertProtobufBidsToDomain(pbAuctionState.GetBidList()),
		SortedTxList: ConvertProtobufTxsToDomain(pbAuctionState.GetSortedTxList()),
		IsEnded:      pbAuctionState.GetIsEnded(),
	}
}

func ConvertDomainAuctionStateToProtobuf(domainAuctionState AuctionState) *auctionpb.AuctionState {
	return &auctionpb.AuctionState{
		AuctionInfo:  ConvertDomainAuctionInfoToProtobuf(domainAuctionState.AuctionInfo),
		BidList:      ConvertDomainBidsToProtobuf(domainAuctionState.BidList),
		SortedTxList: ConvertDomainTxsToProtobuf(domainAuctionState.SortedTxList),
		IsEnded:      domainAuctionState.IsEnded,
	}
}
