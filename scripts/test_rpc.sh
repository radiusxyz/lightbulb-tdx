#!/bin/zsh

# gRPC server configuration
GRPC_URL="localhost:50051"

# Functions for each RPC method
add_auction() {
    # Current timestamp in seconds
    NOW=$(date +%s)

    # Start time: now + 5 seconds
    START=$(($NOW + 5))

    # End time: start + 15 seconds
    END=$(($START + 15))

    # Convert timestamps to milliseconds
    START_MS=$(($START * 1000))
    END_MS=$(($END * 1000))

    # Auction ID (can be passed as an argument or hardcoded)
    AUCTION_ID=${1:-"auction123"}

    # gRPC request payload
    JSON_PAYLOAD=$(
        cat <<EOF
{
  "auction_info": {
    "auction_id": "$AUCTION_ID",
    "chain_id": 1,
    "start_time": $START_MS,
    "end_time": $END_MS,
    "seller_address": "0xSellerAddress",
    "block_number": 1000,
    "blockspace_size": 500,
    "seller_signature": "signature123"
  }
}
EOF
    )

    # Print payload for debugging
    echo "Sending gRPC request to AddAuction with payload:"
    echo "$JSON_PAYLOAD"

    # Execute gRPC call
    grpcurl -plaintext -d "$JSON_PAYLOAD" $GRPC_URL auction.AuctionService/AddAuction
}

submit_bids() {
    CHAIN_ID=1
    AUCTION_ID="auction123"

    # gRPC request payload
    JSON_PAYLOAD=$(
        cat <<EOF
{
  "chain_id": $CHAIN_ID,
  "auction_id": "$AUCTION_ID",
  "bid_list": [
    {
      "auction_id": "$AUCTION_ID",
      "chain_id": $CHAIN_ID,
      "bidder_addr": "0xBidderAddress1",
      "bid_amount": 500,
      "bidder_signature": "signature123",
      "tx_list": [
        {"tx_data": "Transaction1Data"},
        {"tx_data": "Transaction2Data"}
      ]
    },
    {
      "auction_id": "$AUCTION_ID",
      "chain_id": $CHAIN_ID,
      "bidder_addr": "0xBidderAddress2",
      "bid_amount": 600,
      "bidder_signature": "signature456",
      "tx_list": [
        {"tx_data": "Transaction3Data"},
        {"tx_data": "Transaction4Data"}
      ]
    },
    {
      "auction_id": "$AUCTION_ID",
      "chain_id": $CHAIN_ID,
      "bidder_addr": "0xBidderAddress3",
      "bid_amount": 700,
      "bidder_signature": "signature789",
      "tx_list": [
        {"tx_data": "Transaction5Data"},
        {"tx_data": "Transaction6Data"}
      ]
    }
  ]
}
EOF
    )

    # Print payload for debugging
    echo "Sending gRPC request to SubmitBids with payload:"
    echo "$JSON_PAYLOAD"

    # Execute gRPC call
    grpcurl -plaintext -d "$JSON_PAYLOAD" $GRPC_URL auction.AuctionService/SubmitBids
}

get_auction_info() {
    CHAIN_ID=1
    AUCTION_ID="auction123"

    # gRPC request payload
    JSON_PAYLOAD=$(
        cat <<EOF
{
  "chain_id": $CHAIN_ID,
  "auction_id": "$AUCTION_ID"
}
EOF
    )

    # Print payload for debugging
    echo "Sending gRPC request to GetAuctionInfo with payload:"
    echo "$JSON_PAYLOAD"

    # Execute gRPC call
    grpcurl -plaintext -d "$JSON_PAYLOAD" $GRPC_URL auction.AuctionService/GetAuctionInfo
}

get_auction_state() {
    CHAIN_ID=1

    # gRPC request payload
    JSON_PAYLOAD=$(
        cat <<EOF
{
  "chain_id": $CHAIN_ID
}
EOF
    )

    # Print payload for debugging
    echo "Sending gRPC request to GetAuctionState with payload:"
    echo "$JSON_PAYLOAD"

    # Execute gRPC call
    grpcurl -plaintext -d "$JSON_PAYLOAD" $GRPC_URL auction.AuctionService/GetAuctionState
}

get_latest_tob() {
    CHAIN_ID=1

    # gRPC request payload
    JSON_PAYLOAD=$(
        cat <<EOF
{
  "chain_id": $CHAIN_ID
}
EOF
    )

    # Print payload for debugging
    echo "Sending gRPC request to GetLatestTob with payload:"
    echo "$JSON_PAYLOAD"

    # Execute gRPC call
    grpcurl -plaintext -d "$JSON_PAYLOAD" $GRPC_URL auction.AuctionService/GetLatestTob
}

# Main script logic
case $1 in
add_auction)
    add_auction "$2"
    ;;
submit_bids)
    submit_bids
    ;;
get_auction_info)
    get_auction_info
    ;;
get_auction_state)
    get_auction_state
    ;;
get_latest_tob)
    get_latest_tob
    ;;
*)
    echo "Usage: $0 {add_auction|submit_bids|get_auction_info|get_auction_state|get_latest_tob} [arguments]"
    exit 1
    ;;
esac
