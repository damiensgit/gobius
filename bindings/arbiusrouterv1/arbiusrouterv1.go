// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arbiusrouterv1

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// Signature is an auto generated low-level Go binding around an user-defined struct.
type Signature struct {
	Signer    common.Address
	Signature []byte
}

// ArbiusRouterV1MetaData contains all meta data concerning the ArbiusRouterV1 contract.
var ArbiusRouterV1MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"engine_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"arbius_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"router_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InsufficientSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValidator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignersNotSorted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SubmitTaskFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TimeNotPassed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"taskid\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"IncentiveAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"taskid\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"IncentiveClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minValidators\",\"type\":\"uint256\"}],\"name\":\"MinValidatorsSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"ValidatorSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"}],\"name\":\"addIncentive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"arbius\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"taskids_\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structSignature[]\",\"name\":\"sigs_\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"sigsPerTask_\",\"type\":\"uint256\"}],\"name\":\"bulkClaimIncentive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structSignature[]\",\"name\":\"sigs_\",\"type\":\"tuple[]\"}],\"name\":\"claimIncentive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"taskid_\",\"type\":\"bytes32\"}],\"name\":\"emergencyClaimIncentive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"engine\",\"outputs\":[{\"internalType\":\"contractIArbius\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"incentives\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"receiver\",\"outputs\":[{\"internalType\":\"contractSwapReceiver\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router\",\"outputs\":[{\"internalType\":\"contractIUniswapV2Router02\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minValidators_\",\"type\":\"uint256\"}],\"name\":\"setMinValidators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator_\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"status_\",\"type\":\"bool\"}],\"name\":\"setValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"input_\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"incentive_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gas_\",\"type\":\"uint256\"}],\"name\":\"submitTask\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"input_\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"incentive_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gas_\",\"type\":\"uint256\"}],\"name\":\"submitTaskWithETH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"version_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"model_\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"input_\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"incentive_\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"token_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountInMax_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gas_\",\"type\":\"uint256\"}],\"name\":\"submitTaskWithToken\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token_\",\"type\":\"address\"}],\"name\":\"uniswapApprove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash_\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structSignature[]\",\"name\":\"sigs_\",\"type\":\"tuple[]\"}],\"name\":\"validateSignatures\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token_\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawETH\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060409080825234620001c957608081620021fe8038038091620000258285620001ce565b833981010312620001c9576200003b8162000208565b906020906200004c82820162000208565b6200006760606200005f87850162000208565b930162000208565b600080546001600160a01b0319808216339081178455895193989495946001600160a01b03949385928392839283929083167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08f80a31696878560015416176001551697888460025416176002551695868360035416176003551690600454161760045563095ea7b360e01b92838252600482015285816044818a600019998a60248401525af18015620001bf576044928795949289926200019d575b5060025416895196879586948552600485015260248401525af1801562000193576200015e575b505060065551611fc69081620002388239f35b816200018292903d106200018b575b620001798183620001ce565b8101906200021d565b5038806200014b565b503d6200016d565b84513d85823e3d90fd5b620001b790873d89116200018b57620001798183620001ce565b503862000124565b88513d89823e3d90fd5b600080fd5b601f909101601f19168101906001600160401b03821190821017620001f257604052565b634e487b7160e01b600052604160045260246000fd5b51906001600160a01b0382168203620001c957565b90816020910312620001c957518015158103620001c9579056fe6080604052600436101561001b575b361561001957600080fd5b005b6000803560e01c80631cec43ba1461142b578063297035b3146113d15780634623c91d1461135257806351cff8d91461125b5780636077779514611231578063715018a6146111e9578063739edcbe14610c8457806379dfe40c14610c365780638da5cb5b14610c0f57806391067f9014610bb1578063b235d02a14610779578063b28678051461075e578063c5ab224114610740578063c9d4623f14610717578063d7f332b61461063a578063e086e5ec1461060c578063e0d61b07146105e3578063e93ae81c1461045f578063f2fde38b146103ae578063f7260d3e14610385578063f887ea401461035c578063fa52c7d81461031d5763fb53f5b114610124575061000e565b3461031a5780806101343661179c565b95929098969493602061017360018060a01b03600254166101558d8561193e565b6040519b8c809481936323b872dd60e01b8352303360048501611961565b03925af197881561030f5789986102b9575b506101be92889594926101b092604051968795602087019a6308745dd160e01b8c5260248801611983565b03601f1981018352826118a8565b60018060a01b0360015416905193f16101d561190e565b50156102a7576001546040516360beed9560e11b81529290602090849060049082906001600160a01b03165afa92831561029a578193610260575b5091816020938293610228575b505050604051908152f35b600080516020611f318339815191529160408285889452600784522061024f82825461193e565b9055604051908152a280388061021d565b92506020833d602011610292575b8161027b602093836118a8565b8101031261028d576020925192610210565b600080fd5b3d915061026e565b50604051903d90823e3d90fd5b6040516335541d8d60e21b8152600490fd5b90929750602094939194813d602011610307575b816102da602093836118a8565b8101031261030357889788956101be946102f66101b0946118cb565b5092509294955092610185565b8880fd5b3d91506102cd565b6040513d8b823e3d90fd5b80fd5b503461031a57602036600319011261031a5760209060ff906040906001600160a01b03610348611759565b168152600584522054166040519015158152f35b503461031a578060031936011261031a576003546040516001600160a01b039091168152602090f35b503461031a578060031936011261031a576004546040516001600160a01b039091168152602090f35b503461031a57602036600319011261031a576103c8611759565b6103d0611d17565b6001600160a01b0390811690811561040b576000548260018060a01b031982161760005516600080516020611f71833981519152600080a380f35b60405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608490fd5b503461031a578061046f3661182c565b6001546040516375c7050960e01b60208083019182526024808401889052835295946001600160a01b03949390928892839287169183906104af8161188d565b51925af150826104ce6104c061190e565b878082518301019101611c42565b93919290501633036105a7575b5061051f9286949282866104f29451910120611b54565b6002541683865260078352604086205460405196878094819363a9059cbb60e01b835233600484016118d8565b03925af1801561059c57610563575b60079250808452828252604084205460405190815281600080516020611f51833981519152843393a383525280604081205580f35b8183813d8311610595575b61057881836118a8565b810103126105915761058b6007936118cb565b5061052e565b8380fd5b503d61056e565b6040513d86823e3d90fd5b909650919290916001600160401b03906105c090611cef565b1642106105d15791908695386104db565b604051631a0a745560e11b8152600490fd5b503461031a578060031936011261031a576002546040516001600160a01b039091168152602090f35b503461031a578060031936011261031a57610625611d17565b8080808047335af15061063661190e565b5080f35b503461031a57604036600319011261031a5760043560243560018060a01b03600254166040519182916323b872dd60e01b8352602092839181888161068487303360048501611961565b03925af1801561070c576106c6575b600080516020611f31833981519152925083855260078252604085206106ba82825461193e565b9055604051908152a280f35b8183813d8311610705575b6106db81836118a8565b81010312610701576106fb600080516020611f31833981519152936118cb565b50610693565b8480fd5b503d6106d1565b6040513d87823e3d90fd5b503461031a578060031936011261031a576001546040516001600160a01b039091168152602090f35b503461031a578060031936011261031a576020600654604051908152f35b503461031a576107766107703661182c565b91611b54565b80f35b503461031a576060908160031936011261031a576004356001600160401b03811161088e576107ac9036906004016117fc565b906024356001600160401b038111610591576107cc9036906004016117fc565b949060065460443510610b9f576107e76044969496356119cf565b916107f560405193846118a8565b6044358352601f196108086044356119cf565b0190865b828110610b7657505050849585905b8082106108a557868060206108538b60018060a01b03600254169060405194858094819363a9059cbb60e01b835233600484016118d8565b03925af1801561089a57610865575080f35b6020813d8211610892575b8161087d602093836118a8565b8101031261088e57610636906118cb565b5080fd5b3d9150610870565b6040513d84823e3d90fd5b868060018060a09b9596989a97999b1b03600154166108c586858c611d07565b358260405160208101926375c7050960e01b84526024820152602481526108eb8161188d565b51925af15061090a6108fb61190e565b60208082518301019101611c42565b926001600160a01b031633039050610b57575b506020815191012094845b6044358110610a8757508490815b88518310156109f2576001600160a01b03610951848b611a03565b515116906001600160a01b03168111156109e057808752600560205260ff604088205416156109ce578061099e610996602061098d878e611a03565b5101518b611e84565b919091611d6f565b6001600160a01b0316036109bc576109b69092611aec565b91610936565b604051638baa579f60e01b8152600490fd5b604051631a0a9b9f60e21b8152600490fd5b60405163a7781cbb60e01b8152600490fd5b50939050610a80919694979550979197610a0d818888611d07565b358852610a256007998a60205260408a20549061193e565b98610a31828989611d07565b35610a3d838a8a611d07565b358a528160205260408a205490604051918252600080516020611f5183398151915260203393a3610a6f828989611d07565b358952602052876040812055611aec565b909461081b565b60449992993585028581046044351486151715610b4357610aab82610ab29261193e565b8484611afb565b604081360312610b3f5760405190610ac982611872565b80356001600160a01b038116810361028d5782526020810135906001600160401b038211610303570136601f82011215610b3b5790610b13610b3393923690602081359101611b1d565b6020820152610b22828b611a03565b52610b2d818a611a03565b50611aec565b989198610928565b8780fd5b8680fd5b634e487b7160e01b87526011600452602487fd5b6001600160401b0390610b6990611cef565b1642106105d1573861091d565b602090604099979951610b8881611872565b89815282848183015282880101520197959761080c565b604051633724e34360e11b8152600490fd5b503461031a57602036600319011261031a578060206001600160a01b03604481610bd9611759565b60035460405163095ea7b360e01b8152941660048501526000196024850152929485938492165af1801561089a57610865575080f35b503461031a578060031936011261031a57546040516001600160a01b039091168152602090f35b503461031a57602036600319011261031a577fbcfae85be40ac3606c557faf143ce6b08c7d99137b0c98eff034fddc6926c31b6020600435610c76611d17565b80600655604051908152a180f35b503461031a5761012036600319011261031a5760ff600435166004350361028d576024356001600160a01b038116900361028d576001600160401b036084358181116110e057610cd890369060040161176f565b60c43592916001600160a01b038416840361028d576040516323b872dd60e01b815260208180610d0f60e435303360048501611961565b0381896001600160a01b038a165af1801561111d576111b0575b5060405192608084019081118482101761119a576040526003835260603660208501376001600160a01b038416610d5f846119e6565b526003546040516315ab88c960e31b81526001600160a01b039091169390602081600481885afa90811561118f578791611160575b50610d9e826119f3565b6001600160a01b03918216905260028054835192169591111561114a57610e1687959186926060850152610dd660a43560643561193e565b938360018060a01b03600454169560405196879586948593634401edf760e11b8552600485015260e435602485015260a0604485015260a4840190611aaf565b906064830152600019608483015203925af1801561059c57611128575b50600480546002546040516370a0823160e01b81526001600160a01b03928316938101849052929492911690602081602481855afa90811561111d5786916110e8575b50843b156110e457610ea194869283604051809881958294635705ae4360e01b8452600484016118d8565b03925af192831561059c5784936110c6575b5090610ef4610ee69260405193849160208301946308745dd160e01b865260643560443560243560043560248801611983565b03601f1981018452836118a8565b60015491519183906001600160a01b031661010435f1610f1261190e565b50156102a7576001546040516360beed9560e11b81529190602090839060049082906001600160a01b03165afa9182156110bb578392611087575b5060a43561104c575b6040516370a0823160e01b81523060048201526020816024816001600160a01b0386165afa90811561059c57849161101a575b5080610f9b575b602083604051908152f35b610fc0849260209260405196878094819363a9059cbb60e01b835233600484016118d8565b03926001600160a01b03165af1801561029a57610fde575b80610f90565b6020833d602011611012575b81610ff7602093836118a8565b8101031261031a575061100b6020926118cb565b5038610fd8565b3d9150610fea565b90506020813d602011611044575b81611035602093836118a8565b8101031261028d575138610f89565b3d9150611028565b81835260076020526040832061106560a435825461193e565b905581600080516020611f31833981519152602060405160a4358152a2610f56565b9091506020813d6020116110b3575b816110a3602093836118a8565b8101031261028d57519038610f4d565b3d9150611096565b6040513d85823e3d90fd5b6110d3909391929361185f565b6110e05790829138610eb3565b8280fd5b8580fd5b9550506020853d602011611115575b81611104602093836118a8565b8101031261028d5786945138610e76565b3d91506110f7565b6040513d88823e3d90fd5b611143903d8086833e61113b81836118a8565b810190611a36565b5038610e33565b634e487b7160e01b600052603260045260246000fd5b611182915060203d602011611188575b61117a81836118a8565b810190611a17565b38610d94565b503d611170565b6040513d89823e3d90fd5b634e487b7160e01b600052604160045260246000fd5b6020813d6020116111e1575b816111c9602093836118a8565b810103126110e4576111da906118cb565b5038610d29565b3d91506111bc565b503461031a578060031936011261031a57611202611d17565b600080546001600160a01b0319811682556001600160a01b0316600080516020611f718339815191528280a380f35b503461031a57602036600319011261031a5760406020916004358152600783522054604051908152f35b503461031a5760208060031936011261088e57611276611759565b9061127f611d17565b6040516370a0823160e01b81523060048201526001600160a01b03928316928282602481875afa91821561070c578592611321575b50906112dc9383928654168660405180978195829463a9059cbb60e01b8452600484016118d8565b03925af180156110bb576112ee578280f35b81813d831161131a575b61130281836118a8565b8101031261088e57611313906118cb565b5038808280f35b503d6112f8565b91508282813d831161134b575b61133881836118a8565b8101031261028d579051906112dc6112b4565b503d61132e565b503461031a57604036600319011261031a5761136c611759565b602435908115158092036110e05760207f6b3b7d0d26ec99d080840dca1323c7147d1868adc66a4290afb8101d7908320d916113a6611d17565b6001600160a01b0316808552600582526040808620805460ff191660ff87161790555193845292a280f35b503461031a5780602080600319360112611428576004356113f0611d17565b8161051f60018060a01b036002541683865260078352604086205460405196878094819363a9059cbb60e01b835233600484016118d8565b50fd5b506114353661179c565b949261144898969291979498608061188d565b600260805260403660a0376003546040516315ab88c960e31b81526001600160a01b0390911690602081600481855afa90811561030f57899161173a575b5061149160806119e6565b6001600160a01b039182169052600254166114ac60806119f3565b52876114b88a8c61193e565b60018060a01b036004541692604051808095819463fb3bdb4160e01b83526004830152608060248301526114f0608483016080611aaf565b9060448301526000196064830152039134905af180156116d257611720575b50600480546002546040516370a0823160e01b81526001600160a01b039283169381018490529b92911660208c602481845afa9b8c15611715578a9c6116e1575b50813b156116dd579a899161157d9b9c83604051809e81958294635705ae4360e01b8452600484016118d8565b03925af180156116d2576116af575b879850926101b088979695936115bc938996604051968795602087019a6308745dd160e01b8c5260248801611983565b60018060a01b0360015416905193f16115d361190e565b50156102a7576001546040516360beed9560e11b81529190602090839060049082906001600160a01b03165afa91821561029a578192611678575b508183602094611641575b50504761162a575b50604051908152f35b80808047335af15061163a61190e565b5038611621565b84600080516020611f3183398151915291838552600782526040852061166882825461193e565b9055604051908152a28138611619565b9291506020833d6020116116a7575b81611694602093836118a8565b8101031261028d5760209251919261160e565b3d9150611687565b926115bc9196959492976116c56101b09a61185f565b979294959691509261158c565b6040513d8a823e3d90fd5b8980fd5b909b506020813d60201161170d575b816116fd602093836118a8565b810103126116dd57519a38611550565b3d91506116f0565b6040513d8c823e3d90fd5b611733903d808a833e61113b81836118a8565b503861150f565b611753915060203d6020116111885761117a81836118a8565b38611486565b600435906001600160a01b038216820361028d57565b9181601f8401121561028d578235916001600160401b03831161028d576020838186019501011161028d57565b60e060031982011261028d5760043560ff8116810361028d57916024356001600160a01b038116810361028d57916044359160643591608435906001600160401b03821161028d576117f09160040161176f565b909160a4359060c43590565b9181601f8401121561028d578235916001600160401b03831161028d576020808501948460051b01011161028d57565b90604060031983011261028d5760043591602435906001600160401b03821161028d5761185b916004016117fc565b9091565b6001600160401b03811161119a57604052565b604081019081106001600160401b0382111761119a57604052565b606081019081106001600160401b0382111761119a57604052565b601f909101601f19168101906001600160401b0382119082101761119a57604052565b5190811515820361028d57565b6001600160a01b039091168152602081019190915260400190565b6001600160401b03811161119a57601f01601f191660200190565b3d15611939573d9061191f826118f3565b9161192d60405193846118a8565b82523d6000602084013e565b606090565b9190820180921161194b57565b634e487b7160e01b600052601160045260246000fd5b6001600160a01b03918216815291166020820152604081019190915260600190565b9491928694919360ff60c0989516875260018060a01b031660208701526040860152606085015260a060808501528160a0850152848401376000828201840152601f01601f1916010190565b6001600160401b03811161119a5760051b60200190565b80511561114a5760200190565b80516001101561114a5760400190565b805182101561114a5760209160051b010190565b9081602091031261028d57516001600160a01b038116810361028d5790565b602090818184031261028d578051906001600160401b03821161028d57019180601f8401121561028d578251611a6b816119cf565b93611a7960405195866118a8565b818552838086019260051b82010192831161028d578301905b828210611aa0575050505090565b81518152908301908301611a92565b90815180825260208080930193019160005b828110611acf575050505090565b83516001600160a01b031685529381019392810192600101611ac1565b600019811461194b5760010190565b919081101561114a5760051b81013590603e198136030182121561028d570190565b929192611b29826118f3565b91611b3760405193846118a8565b82948184528183011161028d578281602093846000960137010152565b90916006548110610b9f576000805b828210611b71575050505050565b611b7c828487611afb565b35906001600160a01b03908183169081840361028d578282911610156109e057806000526020916005835260409260ff84600020541615611c3157611bc286888b611afb565b8181013590601e198136030182121561028d5701803591906001600160401b03831161028d5701813603811361028d57611c0461099691611c0a933691611b1d565b89611e84565b1603611c215750611c1b9091611aec565b90611b63565b51638baa579f60e01b8152600490fd5b8351631a0a9b9f60e21b8152600490fd5b60808183031261028d5780516001600160a01b038116810361028d5792602092838301519360018060401b0394858116810361028d5794611c85604086016118cb565b94606081015191821161028d570182601f8201121561028d57805190611caa826118f3565b93611cb860405195866118a8565b82855283838301011161028d5760005b828110611cdc575050906000918301015290565b8181018401518582018501528301611cc8565b6001600160401b03908116603c019190821161194b57565b919081101561114a5760051b0190565b6000546001600160a01b03163303611d2b57565b606460405162461bcd60e51b815260206004820152602060248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152fd5b6005811015611e6e5780611d805750565b60018103611dc85760405162461bcd60e51b815260206004820152601860248201527745434453413a20696e76616c6964207369676e617475726560401b6044820152606490fd5b60028103611e155760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e677468006044820152606490fd5b600314611e1e57565b60405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b6064820152608490fd5b634e487b7160e01b600052602160045260246000fd5b906041815114600014611eae5761185b916020820151906060604084015193015160001a90611eb8565b5050600090600290565b9291906fa2a8918ca85bafe22016d0b997e4df60600160ff1b038311611f245791608094939160ff602094604051948552168484015260408301526060820152600093849182805260015afa1561029a5781516001600160a01b03811615611f1e579190565b50600190565b5050505060009060039056fe19d4213c1f22deb153156be5bf9eee8fe77c36a557d83434b8cbb543fc034d6a8c8fb16c3fff3e9353f4a39b33dd9e38ba88594f8c66defbff2048265738780b8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0a2646970667358221220b3d68d4aeea3cd91c6afc3152a9c4ccf11e822a08d022a9ad154608288e2440164736f6c63430008130033",
}

// ArbiusRouterV1ABI is the input ABI used to generate the binding from.
// Deprecated: Use ArbiusRouterV1MetaData.ABI instead.
var ArbiusRouterV1ABI = ArbiusRouterV1MetaData.ABI

// ArbiusRouterV1Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ArbiusRouterV1MetaData.Bin instead.
var ArbiusRouterV1Bin = ArbiusRouterV1MetaData.Bin

// DeployArbiusRouterV1 deploys a new Ethereum contract, binding an instance of ArbiusRouterV1 to it.
func DeployArbiusRouterV1(auth *bind.TransactOpts, backend bind.ContractBackend, engine_ common.Address, arbius_ common.Address, router_ common.Address, receiver_ common.Address) (common.Address, *types.Transaction, *ArbiusRouterV1, error) {
	parsed, err := ArbiusRouterV1MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ArbiusRouterV1Bin), backend, engine_, arbius_, router_, receiver_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ArbiusRouterV1{ArbiusRouterV1Caller: ArbiusRouterV1Caller{contract: contract}, ArbiusRouterV1Transactor: ArbiusRouterV1Transactor{contract: contract}, ArbiusRouterV1Filterer: ArbiusRouterV1Filterer{contract: contract}}, nil
}

// ArbiusRouterV1 is an auto generated Go binding around an Ethereum contract.
type ArbiusRouterV1 struct {
	ArbiusRouterV1Caller     // Read-only binding to the contract
	ArbiusRouterV1Transactor // Write-only binding to the contract
	ArbiusRouterV1Filterer   // Log filterer for contract events
}

// ArbiusRouterV1Caller is an auto generated read-only Go binding around an Ethereum contract.
type ArbiusRouterV1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArbiusRouterV1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ArbiusRouterV1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArbiusRouterV1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ArbiusRouterV1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArbiusRouterV1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ArbiusRouterV1Session struct {
	Contract     *ArbiusRouterV1   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ArbiusRouterV1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ArbiusRouterV1CallerSession struct {
	Contract *ArbiusRouterV1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ArbiusRouterV1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ArbiusRouterV1TransactorSession struct {
	Contract     *ArbiusRouterV1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ArbiusRouterV1Raw is an auto generated low-level Go binding around an Ethereum contract.
type ArbiusRouterV1Raw struct {
	Contract *ArbiusRouterV1 // Generic contract binding to access the raw methods on
}

// ArbiusRouterV1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ArbiusRouterV1CallerRaw struct {
	Contract *ArbiusRouterV1Caller // Generic read-only contract binding to access the raw methods on
}

// ArbiusRouterV1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ArbiusRouterV1TransactorRaw struct {
	Contract *ArbiusRouterV1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewArbiusRouterV1 creates a new instance of ArbiusRouterV1, bound to a specific deployed contract.
func NewArbiusRouterV1(address common.Address, backend bind.ContractBackend) (*ArbiusRouterV1, error) {
	contract, err := bindArbiusRouterV1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArbiusRouterV1{ArbiusRouterV1Caller: ArbiusRouterV1Caller{contract: contract}, ArbiusRouterV1Transactor: ArbiusRouterV1Transactor{contract: contract}, ArbiusRouterV1Filterer: ArbiusRouterV1Filterer{contract: contract}}, nil
}

// NewArbiusRouterV1Caller creates a new read-only instance of ArbiusRouterV1, bound to a specific deployed contract.
func NewArbiusRouterV1Caller(address common.Address, caller bind.ContractCaller) (*ArbiusRouterV1Caller, error) {
	contract, err := bindArbiusRouterV1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArbiusRouterV1Caller{contract: contract}, nil
}

// NewArbiusRouterV1Transactor creates a new write-only instance of ArbiusRouterV1, bound to a specific deployed contract.
func NewArbiusRouterV1Transactor(address common.Address, transactor bind.ContractTransactor) (*ArbiusRouterV1Transactor, error) {
	contract, err := bindArbiusRouterV1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArbiusRouterV1Transactor{contract: contract}, nil
}

// NewArbiusRouterV1Filterer creates a new log filterer instance of ArbiusRouterV1, bound to a specific deployed contract.
func NewArbiusRouterV1Filterer(address common.Address, filterer bind.ContractFilterer) (*ArbiusRouterV1Filterer, error) {
	contract, err := bindArbiusRouterV1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArbiusRouterV1Filterer{contract: contract}, nil
}

// bindArbiusRouterV1 binds a generic wrapper to an already deployed contract.
func bindArbiusRouterV1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ArbiusRouterV1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ArbiusRouterV1 *ArbiusRouterV1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbiusRouterV1.Contract.ArbiusRouterV1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ArbiusRouterV1 *ArbiusRouterV1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.ArbiusRouterV1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ArbiusRouterV1 *ArbiusRouterV1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.ArbiusRouterV1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ArbiusRouterV1 *ArbiusRouterV1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArbiusRouterV1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.contract.Transact(opts, method, params...)
}

// Arbius is a free data retrieval call binding the contract method 0xe0d61b07.
//
// Solidity: function arbius() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1Caller) Arbius(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbiusRouterV1.contract.Call(opts, &out, "arbius")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Arbius is a free data retrieval call binding the contract method 0xe0d61b07.
//
// Solidity: function arbius() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) Arbius() (common.Address, error) {
	return _ArbiusRouterV1.Contract.Arbius(&_ArbiusRouterV1.CallOpts)
}

// Arbius is a free data retrieval call binding the contract method 0xe0d61b07.
//
// Solidity: function arbius() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1CallerSession) Arbius() (common.Address, error) {
	return _ArbiusRouterV1.Contract.Arbius(&_ArbiusRouterV1.CallOpts)
}

// Engine is a free data retrieval call binding the contract method 0xc9d4623f.
//
// Solidity: function engine() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1Caller) Engine(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbiusRouterV1.contract.Call(opts, &out, "engine")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Engine is a free data retrieval call binding the contract method 0xc9d4623f.
//
// Solidity: function engine() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) Engine() (common.Address, error) {
	return _ArbiusRouterV1.Contract.Engine(&_ArbiusRouterV1.CallOpts)
}

// Engine is a free data retrieval call binding the contract method 0xc9d4623f.
//
// Solidity: function engine() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1CallerSession) Engine() (common.Address, error) {
	return _ArbiusRouterV1.Contract.Engine(&_ArbiusRouterV1.CallOpts)
}

// Incentives is a free data retrieval call binding the contract method 0x60777795.
//
// Solidity: function incentives(bytes32 ) view returns(uint256)
func (_ArbiusRouterV1 *ArbiusRouterV1Caller) Incentives(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _ArbiusRouterV1.contract.Call(opts, &out, "incentives", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Incentives is a free data retrieval call binding the contract method 0x60777795.
//
// Solidity: function incentives(bytes32 ) view returns(uint256)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) Incentives(arg0 [32]byte) (*big.Int, error) {
	return _ArbiusRouterV1.Contract.Incentives(&_ArbiusRouterV1.CallOpts, arg0)
}

// Incentives is a free data retrieval call binding the contract method 0x60777795.
//
// Solidity: function incentives(bytes32 ) view returns(uint256)
func (_ArbiusRouterV1 *ArbiusRouterV1CallerSession) Incentives(arg0 [32]byte) (*big.Int, error) {
	return _ArbiusRouterV1.Contract.Incentives(&_ArbiusRouterV1.CallOpts, arg0)
}

// MinValidators is a free data retrieval call binding the contract method 0xc5ab2241.
//
// Solidity: function minValidators() view returns(uint256)
func (_ArbiusRouterV1 *ArbiusRouterV1Caller) MinValidators(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArbiusRouterV1.contract.Call(opts, &out, "minValidators")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinValidators is a free data retrieval call binding the contract method 0xc5ab2241.
//
// Solidity: function minValidators() view returns(uint256)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) MinValidators() (*big.Int, error) {
	return _ArbiusRouterV1.Contract.MinValidators(&_ArbiusRouterV1.CallOpts)
}

// MinValidators is a free data retrieval call binding the contract method 0xc5ab2241.
//
// Solidity: function minValidators() view returns(uint256)
func (_ArbiusRouterV1 *ArbiusRouterV1CallerSession) MinValidators() (*big.Int, error) {
	return _ArbiusRouterV1.Contract.MinValidators(&_ArbiusRouterV1.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbiusRouterV1.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) Owner() (common.Address, error) {
	return _ArbiusRouterV1.Contract.Owner(&_ArbiusRouterV1.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1CallerSession) Owner() (common.Address, error) {
	return _ArbiusRouterV1.Contract.Owner(&_ArbiusRouterV1.CallOpts)
}

// Receiver is a free data retrieval call binding the contract method 0xf7260d3e.
//
// Solidity: function receiver() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1Caller) Receiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbiusRouterV1.contract.Call(opts, &out, "receiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Receiver is a free data retrieval call binding the contract method 0xf7260d3e.
//
// Solidity: function receiver() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) Receiver() (common.Address, error) {
	return _ArbiusRouterV1.Contract.Receiver(&_ArbiusRouterV1.CallOpts)
}

// Receiver is a free data retrieval call binding the contract method 0xf7260d3e.
//
// Solidity: function receiver() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1CallerSession) Receiver() (common.Address, error) {
	return _ArbiusRouterV1.Contract.Receiver(&_ArbiusRouterV1.CallOpts)
}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1Caller) Router(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArbiusRouterV1.contract.Call(opts, &out, "router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) Router() (common.Address, error) {
	return _ArbiusRouterV1.Contract.Router(&_ArbiusRouterV1.CallOpts)
}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_ArbiusRouterV1 *ArbiusRouterV1CallerSession) Router() (common.Address, error) {
	return _ArbiusRouterV1.Contract.Router(&_ArbiusRouterV1.CallOpts)
}

// ValidateSignatures is a free data retrieval call binding the contract method 0xb2867805.
//
// Solidity: function validateSignatures(bytes32 hash_, (address,bytes)[] sigs_) view returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Caller) ValidateSignatures(opts *bind.CallOpts, hash_ [32]byte, sigs_ []Signature) error {
	var out []interface{}
	err := _ArbiusRouterV1.contract.Call(opts, &out, "validateSignatures", hash_, sigs_)

	if err != nil {
		return err
	}

	return err

}

// ValidateSignatures is a free data retrieval call binding the contract method 0xb2867805.
//
// Solidity: function validateSignatures(bytes32 hash_, (address,bytes)[] sigs_) view returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) ValidateSignatures(hash_ [32]byte, sigs_ []Signature) error {
	return _ArbiusRouterV1.Contract.ValidateSignatures(&_ArbiusRouterV1.CallOpts, hash_, sigs_)
}

// ValidateSignatures is a free data retrieval call binding the contract method 0xb2867805.
//
// Solidity: function validateSignatures(bytes32 hash_, (address,bytes)[] sigs_) view returns()
func (_ArbiusRouterV1 *ArbiusRouterV1CallerSession) ValidateSignatures(hash_ [32]byte, sigs_ []Signature) error {
	return _ArbiusRouterV1.Contract.ValidateSignatures(&_ArbiusRouterV1.CallOpts, hash_, sigs_)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_ArbiusRouterV1 *ArbiusRouterV1Caller) Validators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _ArbiusRouterV1.contract.Call(opts, &out, "validators", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) Validators(arg0 common.Address) (bool, error) {
	return _ArbiusRouterV1.Contract.Validators(&_ArbiusRouterV1.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0xfa52c7d8.
//
// Solidity: function validators(address ) view returns(bool)
func (_ArbiusRouterV1 *ArbiusRouterV1CallerSession) Validators(arg0 common.Address) (bool, error) {
	return _ArbiusRouterV1.Contract.Validators(&_ArbiusRouterV1.CallOpts, arg0)
}

// AddIncentive is a paid mutator transaction binding the contract method 0xd7f332b6.
//
// Solidity: function addIncentive(bytes32 taskid_, uint256 amount_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) AddIncentive(opts *bind.TransactOpts, taskid_ [32]byte, amount_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "addIncentive", taskid_, amount_)
}

// AddIncentive is a paid mutator transaction binding the contract method 0xd7f332b6.
//
// Solidity: function addIncentive(bytes32 taskid_, uint256 amount_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) AddIncentive(taskid_ [32]byte, amount_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.AddIncentive(&_ArbiusRouterV1.TransactOpts, taskid_, amount_)
}

// AddIncentive is a paid mutator transaction binding the contract method 0xd7f332b6.
//
// Solidity: function addIncentive(bytes32 taskid_, uint256 amount_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) AddIncentive(taskid_ [32]byte, amount_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.AddIncentive(&_ArbiusRouterV1.TransactOpts, taskid_, amount_)
}

// BulkClaimIncentive is a paid mutator transaction binding the contract method 0xb235d02a.
//
// Solidity: function bulkClaimIncentive(bytes32[] taskids_, (address,bytes)[] sigs_, uint256 sigsPerTask_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) BulkClaimIncentive(opts *bind.TransactOpts, taskids_ [][32]byte, sigs_ []Signature, sigsPerTask_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "bulkClaimIncentive", taskids_, sigs_, sigsPerTask_)
}

// BulkClaimIncentive is a paid mutator transaction binding the contract method 0xb235d02a.
//
// Solidity: function bulkClaimIncentive(bytes32[] taskids_, (address,bytes)[] sigs_, uint256 sigsPerTask_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) BulkClaimIncentive(taskids_ [][32]byte, sigs_ []Signature, sigsPerTask_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.BulkClaimIncentive(&_ArbiusRouterV1.TransactOpts, taskids_, sigs_, sigsPerTask_)
}

// BulkClaimIncentive is a paid mutator transaction binding the contract method 0xb235d02a.
//
// Solidity: function bulkClaimIncentive(bytes32[] taskids_, (address,bytes)[] sigs_, uint256 sigsPerTask_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) BulkClaimIncentive(taskids_ [][32]byte, sigs_ []Signature, sigsPerTask_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.BulkClaimIncentive(&_ArbiusRouterV1.TransactOpts, taskids_, sigs_, sigsPerTask_)
}

// ClaimIncentive is a paid mutator transaction binding the contract method 0xe93ae81c.
//
// Solidity: function claimIncentive(bytes32 taskid_, (address,bytes)[] sigs_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) ClaimIncentive(opts *bind.TransactOpts, taskid_ [32]byte, sigs_ []Signature) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "claimIncentive", taskid_, sigs_)
}

// ClaimIncentive is a paid mutator transaction binding the contract method 0xe93ae81c.
//
// Solidity: function claimIncentive(bytes32 taskid_, (address,bytes)[] sigs_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) ClaimIncentive(taskid_ [32]byte, sigs_ []Signature) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.ClaimIncentive(&_ArbiusRouterV1.TransactOpts, taskid_, sigs_)
}

// ClaimIncentive is a paid mutator transaction binding the contract method 0xe93ae81c.
//
// Solidity: function claimIncentive(bytes32 taskid_, (address,bytes)[] sigs_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) ClaimIncentive(taskid_ [32]byte, sigs_ []Signature) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.ClaimIncentive(&_ArbiusRouterV1.TransactOpts, taskid_, sigs_)
}

// EmergencyClaimIncentive is a paid mutator transaction binding the contract method 0x297035b3.
//
// Solidity: function emergencyClaimIncentive(bytes32 taskid_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) EmergencyClaimIncentive(opts *bind.TransactOpts, taskid_ [32]byte) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "emergencyClaimIncentive", taskid_)
}

// EmergencyClaimIncentive is a paid mutator transaction binding the contract method 0x297035b3.
//
// Solidity: function emergencyClaimIncentive(bytes32 taskid_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) EmergencyClaimIncentive(taskid_ [32]byte) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.EmergencyClaimIncentive(&_ArbiusRouterV1.TransactOpts, taskid_)
}

// EmergencyClaimIncentive is a paid mutator transaction binding the contract method 0x297035b3.
//
// Solidity: function emergencyClaimIncentive(bytes32 taskid_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) EmergencyClaimIncentive(taskid_ [32]byte) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.EmergencyClaimIncentive(&_ArbiusRouterV1.TransactOpts, taskid_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) RenounceOwnership() (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.RenounceOwnership(&_ArbiusRouterV1.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.RenounceOwnership(&_ArbiusRouterV1.TransactOpts)
}

// SetMinValidators is a paid mutator transaction binding the contract method 0x79dfe40c.
//
// Solidity: function setMinValidators(uint256 minValidators_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) SetMinValidators(opts *bind.TransactOpts, minValidators_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "setMinValidators", minValidators_)
}

// SetMinValidators is a paid mutator transaction binding the contract method 0x79dfe40c.
//
// Solidity: function setMinValidators(uint256 minValidators_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) SetMinValidators(minValidators_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.SetMinValidators(&_ArbiusRouterV1.TransactOpts, minValidators_)
}

// SetMinValidators is a paid mutator transaction binding the contract method 0x79dfe40c.
//
// Solidity: function setMinValidators(uint256 minValidators_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) SetMinValidators(minValidators_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.SetMinValidators(&_ArbiusRouterV1.TransactOpts, minValidators_)
}

// SetValidator is a paid mutator transaction binding the contract method 0x4623c91d.
//
// Solidity: function setValidator(address validator_, bool status_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) SetValidator(opts *bind.TransactOpts, validator_ common.Address, status_ bool) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "setValidator", validator_, status_)
}

// SetValidator is a paid mutator transaction binding the contract method 0x4623c91d.
//
// Solidity: function setValidator(address validator_, bool status_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) SetValidator(validator_ common.Address, status_ bool) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.SetValidator(&_ArbiusRouterV1.TransactOpts, validator_, status_)
}

// SetValidator is a paid mutator transaction binding the contract method 0x4623c91d.
//
// Solidity: function setValidator(address validator_, bool status_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) SetValidator(validator_ common.Address, status_ bool) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.SetValidator(&_ArbiusRouterV1.TransactOpts, validator_, status_)
}

// SubmitTask is a paid mutator transaction binding the contract method 0xfb53f5b1.
//
// Solidity: function submitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 incentive_, uint256 gas_) returns(bytes32)
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) SubmitTask(opts *bind.TransactOpts, version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, incentive_ *big.Int, gas_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "submitTask", version_, owner_, model_, fee_, input_, incentive_, gas_)
}

// SubmitTask is a paid mutator transaction binding the contract method 0xfb53f5b1.
//
// Solidity: function submitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 incentive_, uint256 gas_) returns(bytes32)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) SubmitTask(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, incentive_ *big.Int, gas_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.SubmitTask(&_ArbiusRouterV1.TransactOpts, version_, owner_, model_, fee_, input_, incentive_, gas_)
}

// SubmitTask is a paid mutator transaction binding the contract method 0xfb53f5b1.
//
// Solidity: function submitTask(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 incentive_, uint256 gas_) returns(bytes32)
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) SubmitTask(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, incentive_ *big.Int, gas_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.SubmitTask(&_ArbiusRouterV1.TransactOpts, version_, owner_, model_, fee_, input_, incentive_, gas_)
}

// SubmitTaskWithETH is a paid mutator transaction binding the contract method 0x1cec43ba.
//
// Solidity: function submitTaskWithETH(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 incentive_, uint256 gas_) payable returns(bytes32)
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) SubmitTaskWithETH(opts *bind.TransactOpts, version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, incentive_ *big.Int, gas_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "submitTaskWithETH", version_, owner_, model_, fee_, input_, incentive_, gas_)
}

// SubmitTaskWithETH is a paid mutator transaction binding the contract method 0x1cec43ba.
//
// Solidity: function submitTaskWithETH(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 incentive_, uint256 gas_) payable returns(bytes32)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) SubmitTaskWithETH(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, incentive_ *big.Int, gas_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.SubmitTaskWithETH(&_ArbiusRouterV1.TransactOpts, version_, owner_, model_, fee_, input_, incentive_, gas_)
}

// SubmitTaskWithETH is a paid mutator transaction binding the contract method 0x1cec43ba.
//
// Solidity: function submitTaskWithETH(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 incentive_, uint256 gas_) payable returns(bytes32)
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) SubmitTaskWithETH(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, incentive_ *big.Int, gas_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.SubmitTaskWithETH(&_ArbiusRouterV1.TransactOpts, version_, owner_, model_, fee_, input_, incentive_, gas_)
}

// SubmitTaskWithToken is a paid mutator transaction binding the contract method 0x739edcbe.
//
// Solidity: function submitTaskWithToken(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 incentive_, address token_, uint256 amountInMax_, uint256 gas_) returns(bytes32)
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) SubmitTaskWithToken(opts *bind.TransactOpts, version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, incentive_ *big.Int, token_ common.Address, amountInMax_ *big.Int, gas_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "submitTaskWithToken", version_, owner_, model_, fee_, input_, incentive_, token_, amountInMax_, gas_)
}

// SubmitTaskWithToken is a paid mutator transaction binding the contract method 0x739edcbe.
//
// Solidity: function submitTaskWithToken(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 incentive_, address token_, uint256 amountInMax_, uint256 gas_) returns(bytes32)
func (_ArbiusRouterV1 *ArbiusRouterV1Session) SubmitTaskWithToken(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, incentive_ *big.Int, token_ common.Address, amountInMax_ *big.Int, gas_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.SubmitTaskWithToken(&_ArbiusRouterV1.TransactOpts, version_, owner_, model_, fee_, input_, incentive_, token_, amountInMax_, gas_)
}

// SubmitTaskWithToken is a paid mutator transaction binding the contract method 0x739edcbe.
//
// Solidity: function submitTaskWithToken(uint8 version_, address owner_, bytes32 model_, uint256 fee_, bytes input_, uint256 incentive_, address token_, uint256 amountInMax_, uint256 gas_) returns(bytes32)
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) SubmitTaskWithToken(version_ uint8, owner_ common.Address, model_ [32]byte, fee_ *big.Int, input_ []byte, incentive_ *big.Int, token_ common.Address, amountInMax_ *big.Int, gas_ *big.Int) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.SubmitTaskWithToken(&_ArbiusRouterV1.TransactOpts, version_, owner_, model_, fee_, input_, incentive_, token_, amountInMax_, gas_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.TransferOwnership(&_ArbiusRouterV1.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.TransferOwnership(&_ArbiusRouterV1.TransactOpts, newOwner)
}

// UniswapApprove is a paid mutator transaction binding the contract method 0x91067f90.
//
// Solidity: function uniswapApprove(address token_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) UniswapApprove(opts *bind.TransactOpts, token_ common.Address) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "uniswapApprove", token_)
}

// UniswapApprove is a paid mutator transaction binding the contract method 0x91067f90.
//
// Solidity: function uniswapApprove(address token_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) UniswapApprove(token_ common.Address) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.UniswapApprove(&_ArbiusRouterV1.TransactOpts, token_)
}

// UniswapApprove is a paid mutator transaction binding the contract method 0x91067f90.
//
// Solidity: function uniswapApprove(address token_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) UniswapApprove(token_ common.Address) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.UniswapApprove(&_ArbiusRouterV1.TransactOpts, token_)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address token_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) Withdraw(opts *bind.TransactOpts, token_ common.Address) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "withdraw", token_)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address token_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) Withdraw(token_ common.Address) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.Withdraw(&_ArbiusRouterV1.TransactOpts, token_)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address token_) returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) Withdraw(token_ common.Address) (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.Withdraw(&_ArbiusRouterV1.TransactOpts, token_)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0xe086e5ec.
//
// Solidity: function withdrawETH() returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) WithdrawETH(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.Transact(opts, "withdrawETH")
}

// WithdrawETH is a paid mutator transaction binding the contract method 0xe086e5ec.
//
// Solidity: function withdrawETH() returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) WithdrawETH() (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.WithdrawETH(&_ArbiusRouterV1.TransactOpts)
}

// WithdrawETH is a paid mutator transaction binding the contract method 0xe086e5ec.
//
// Solidity: function withdrawETH() returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) WithdrawETH() (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.WithdrawETH(&_ArbiusRouterV1.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArbiusRouterV1.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ArbiusRouterV1 *ArbiusRouterV1Session) Receive() (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.Receive(&_ArbiusRouterV1.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ArbiusRouterV1 *ArbiusRouterV1TransactorSession) Receive() (*types.Transaction, error) {
	return _ArbiusRouterV1.Contract.Receive(&_ArbiusRouterV1.TransactOpts)
}

// ArbiusRouterV1IncentiveAddedIterator is returned from FilterIncentiveAdded and is used to iterate over the raw logs and unpacked data for IncentiveAdded events raised by the ArbiusRouterV1 contract.
type ArbiusRouterV1IncentiveAddedIterator struct {
	Event *ArbiusRouterV1IncentiveAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ArbiusRouterV1IncentiveAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbiusRouterV1IncentiveAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ArbiusRouterV1IncentiveAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ArbiusRouterV1IncentiveAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArbiusRouterV1IncentiveAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArbiusRouterV1IncentiveAdded represents a IncentiveAdded event raised by the ArbiusRouterV1 contract.
type ArbiusRouterV1IncentiveAdded struct {
	Taskid [32]byte
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterIncentiveAdded is a free log retrieval operation binding the contract event 0x19d4213c1f22deb153156be5bf9eee8fe77c36a557d83434b8cbb543fc034d6a.
//
// Solidity: event IncentiveAdded(bytes32 indexed taskid, uint256 amount)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) FilterIncentiveAdded(opts *bind.FilterOpts, taskid [][32]byte) (*ArbiusRouterV1IncentiveAddedIterator, error) {

	var taskidRule []interface{}
	for _, taskidItem := range taskid {
		taskidRule = append(taskidRule, taskidItem)
	}

	logs, sub, err := _ArbiusRouterV1.contract.FilterLogs(opts, "IncentiveAdded", taskidRule)
	if err != nil {
		return nil, err
	}
	return &ArbiusRouterV1IncentiveAddedIterator{contract: _ArbiusRouterV1.contract, event: "IncentiveAdded", logs: logs, sub: sub}, nil
}

// WatchIncentiveAdded is a free log subscription operation binding the contract event 0x19d4213c1f22deb153156be5bf9eee8fe77c36a557d83434b8cbb543fc034d6a.
//
// Solidity: event IncentiveAdded(bytes32 indexed taskid, uint256 amount)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) WatchIncentiveAdded(opts *bind.WatchOpts, sink chan<- *ArbiusRouterV1IncentiveAdded, taskid [][32]byte) (event.Subscription, error) {

	var taskidRule []interface{}
	for _, taskidItem := range taskid {
		taskidRule = append(taskidRule, taskidItem)
	}

	logs, sub, err := _ArbiusRouterV1.contract.WatchLogs(opts, "IncentiveAdded", taskidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArbiusRouterV1IncentiveAdded)
				if err := _ArbiusRouterV1.contract.UnpackLog(event, "IncentiveAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIncentiveAdded is a log parse operation binding the contract event 0x19d4213c1f22deb153156be5bf9eee8fe77c36a557d83434b8cbb543fc034d6a.
//
// Solidity: event IncentiveAdded(bytes32 indexed taskid, uint256 amount)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) ParseIncentiveAdded(log types.Log) (*ArbiusRouterV1IncentiveAdded, error) {
	event := new(ArbiusRouterV1IncentiveAdded)
	if err := _ArbiusRouterV1.contract.UnpackLog(event, "IncentiveAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ArbiusRouterV1IncentiveClaimedIterator is returned from FilterIncentiveClaimed and is used to iterate over the raw logs and unpacked data for IncentiveClaimed events raised by the ArbiusRouterV1 contract.
type ArbiusRouterV1IncentiveClaimedIterator struct {
	Event *ArbiusRouterV1IncentiveClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ArbiusRouterV1IncentiveClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbiusRouterV1IncentiveClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ArbiusRouterV1IncentiveClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ArbiusRouterV1IncentiveClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArbiusRouterV1IncentiveClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArbiusRouterV1IncentiveClaimed represents a IncentiveClaimed event raised by the ArbiusRouterV1 contract.
type ArbiusRouterV1IncentiveClaimed struct {
	Taskid    [32]byte
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterIncentiveClaimed is a free log retrieval operation binding the contract event 0x8c8fb16c3fff3e9353f4a39b33dd9e38ba88594f8c66defbff2048265738780b.
//
// Solidity: event IncentiveClaimed(bytes32 indexed taskid, address indexed recipient, uint256 amount)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) FilterIncentiveClaimed(opts *bind.FilterOpts, taskid [][32]byte, recipient []common.Address) (*ArbiusRouterV1IncentiveClaimedIterator, error) {

	var taskidRule []interface{}
	for _, taskidItem := range taskid {
		taskidRule = append(taskidRule, taskidItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _ArbiusRouterV1.contract.FilterLogs(opts, "IncentiveClaimed", taskidRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &ArbiusRouterV1IncentiveClaimedIterator{contract: _ArbiusRouterV1.contract, event: "IncentiveClaimed", logs: logs, sub: sub}, nil
}

// WatchIncentiveClaimed is a free log subscription operation binding the contract event 0x8c8fb16c3fff3e9353f4a39b33dd9e38ba88594f8c66defbff2048265738780b.
//
// Solidity: event IncentiveClaimed(bytes32 indexed taskid, address indexed recipient, uint256 amount)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) WatchIncentiveClaimed(opts *bind.WatchOpts, sink chan<- *ArbiusRouterV1IncentiveClaimed, taskid [][32]byte, recipient []common.Address) (event.Subscription, error) {

	var taskidRule []interface{}
	for _, taskidItem := range taskid {
		taskidRule = append(taskidRule, taskidItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _ArbiusRouterV1.contract.WatchLogs(opts, "IncentiveClaimed", taskidRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArbiusRouterV1IncentiveClaimed)
				if err := _ArbiusRouterV1.contract.UnpackLog(event, "IncentiveClaimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIncentiveClaimed is a log parse operation binding the contract event 0x8c8fb16c3fff3e9353f4a39b33dd9e38ba88594f8c66defbff2048265738780b.
//
// Solidity: event IncentiveClaimed(bytes32 indexed taskid, address indexed recipient, uint256 amount)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) ParseIncentiveClaimed(log types.Log) (*ArbiusRouterV1IncentiveClaimed, error) {
	event := new(ArbiusRouterV1IncentiveClaimed)
	if err := _ArbiusRouterV1.contract.UnpackLog(event, "IncentiveClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ArbiusRouterV1MinValidatorsSetIterator is returned from FilterMinValidatorsSet and is used to iterate over the raw logs and unpacked data for MinValidatorsSet events raised by the ArbiusRouterV1 contract.
type ArbiusRouterV1MinValidatorsSetIterator struct {
	Event *ArbiusRouterV1MinValidatorsSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ArbiusRouterV1MinValidatorsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbiusRouterV1MinValidatorsSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ArbiusRouterV1MinValidatorsSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ArbiusRouterV1MinValidatorsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArbiusRouterV1MinValidatorsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArbiusRouterV1MinValidatorsSet represents a MinValidatorsSet event raised by the ArbiusRouterV1 contract.
type ArbiusRouterV1MinValidatorsSet struct {
	MinValidators *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMinValidatorsSet is a free log retrieval operation binding the contract event 0xbcfae85be40ac3606c557faf143ce6b08c7d99137b0c98eff034fddc6926c31b.
//
// Solidity: event MinValidatorsSet(uint256 minValidators)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) FilterMinValidatorsSet(opts *bind.FilterOpts) (*ArbiusRouterV1MinValidatorsSetIterator, error) {

	logs, sub, err := _ArbiusRouterV1.contract.FilterLogs(opts, "MinValidatorsSet")
	if err != nil {
		return nil, err
	}
	return &ArbiusRouterV1MinValidatorsSetIterator{contract: _ArbiusRouterV1.contract, event: "MinValidatorsSet", logs: logs, sub: sub}, nil
}

// WatchMinValidatorsSet is a free log subscription operation binding the contract event 0xbcfae85be40ac3606c557faf143ce6b08c7d99137b0c98eff034fddc6926c31b.
//
// Solidity: event MinValidatorsSet(uint256 minValidators)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) WatchMinValidatorsSet(opts *bind.WatchOpts, sink chan<- *ArbiusRouterV1MinValidatorsSet) (event.Subscription, error) {

	logs, sub, err := _ArbiusRouterV1.contract.WatchLogs(opts, "MinValidatorsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArbiusRouterV1MinValidatorsSet)
				if err := _ArbiusRouterV1.contract.UnpackLog(event, "MinValidatorsSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMinValidatorsSet is a log parse operation binding the contract event 0xbcfae85be40ac3606c557faf143ce6b08c7d99137b0c98eff034fddc6926c31b.
//
// Solidity: event MinValidatorsSet(uint256 minValidators)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) ParseMinValidatorsSet(log types.Log) (*ArbiusRouterV1MinValidatorsSet, error) {
	event := new(ArbiusRouterV1MinValidatorsSet)
	if err := _ArbiusRouterV1.contract.UnpackLog(event, "MinValidatorsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ArbiusRouterV1OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ArbiusRouterV1 contract.
type ArbiusRouterV1OwnershipTransferredIterator struct {
	Event *ArbiusRouterV1OwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ArbiusRouterV1OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbiusRouterV1OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ArbiusRouterV1OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ArbiusRouterV1OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArbiusRouterV1OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArbiusRouterV1OwnershipTransferred represents a OwnershipTransferred event raised by the ArbiusRouterV1 contract.
type ArbiusRouterV1OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ArbiusRouterV1OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ArbiusRouterV1.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ArbiusRouterV1OwnershipTransferredIterator{contract: _ArbiusRouterV1.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ArbiusRouterV1OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ArbiusRouterV1.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArbiusRouterV1OwnershipTransferred)
				if err := _ArbiusRouterV1.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) ParseOwnershipTransferred(log types.Log) (*ArbiusRouterV1OwnershipTransferred, error) {
	event := new(ArbiusRouterV1OwnershipTransferred)
	if err := _ArbiusRouterV1.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ArbiusRouterV1ValidatorSetIterator is returned from FilterValidatorSet and is used to iterate over the raw logs and unpacked data for ValidatorSet events raised by the ArbiusRouterV1 contract.
type ArbiusRouterV1ValidatorSetIterator struct {
	Event *ArbiusRouterV1ValidatorSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ArbiusRouterV1ValidatorSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArbiusRouterV1ValidatorSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ArbiusRouterV1ValidatorSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ArbiusRouterV1ValidatorSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArbiusRouterV1ValidatorSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArbiusRouterV1ValidatorSet represents a ValidatorSet event raised by the ArbiusRouterV1 contract.
type ArbiusRouterV1ValidatorSet struct {
	Validator common.Address
	Status    bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValidatorSet is a free log retrieval operation binding the contract event 0x6b3b7d0d26ec99d080840dca1323c7147d1868adc66a4290afb8101d7908320d.
//
// Solidity: event ValidatorSet(address indexed validator, bool status)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) FilterValidatorSet(opts *bind.FilterOpts, validator []common.Address) (*ArbiusRouterV1ValidatorSetIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _ArbiusRouterV1.contract.FilterLogs(opts, "ValidatorSet", validatorRule)
	if err != nil {
		return nil, err
	}
	return &ArbiusRouterV1ValidatorSetIterator{contract: _ArbiusRouterV1.contract, event: "ValidatorSet", logs: logs, sub: sub}, nil
}

// WatchValidatorSet is a free log subscription operation binding the contract event 0x6b3b7d0d26ec99d080840dca1323c7147d1868adc66a4290afb8101d7908320d.
//
// Solidity: event ValidatorSet(address indexed validator, bool status)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) WatchValidatorSet(opts *bind.WatchOpts, sink chan<- *ArbiusRouterV1ValidatorSet, validator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}

	logs, sub, err := _ArbiusRouterV1.contract.WatchLogs(opts, "ValidatorSet", validatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArbiusRouterV1ValidatorSet)
				if err := _ArbiusRouterV1.contract.UnpackLog(event, "ValidatorSet", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorSet is a log parse operation binding the contract event 0x6b3b7d0d26ec99d080840dca1323c7147d1868adc66a4290afb8101d7908320d.
//
// Solidity: event ValidatorSet(address indexed validator, bool status)
func (_ArbiusRouterV1 *ArbiusRouterV1Filterer) ParseValidatorSet(log types.Log) (*ArbiusRouterV1ValidatorSet, error) {
	event := new(ArbiusRouterV1ValidatorSet)
	if err := _ArbiusRouterV1.contract.UnpackLog(event, "ValidatorSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
