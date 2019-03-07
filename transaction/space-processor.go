package transaction

import (
	"github.com/bottos-project/bottos/role"
	//"github.com/bottos-project/bottos/config"
)

func GetAccountNetLimitt(role role.RoleInterface, sender string)   {

	//sb, _ := role.GetStakedBalance(sender)

	//network_capacityInWindow := config.MAX_SPACE_PER_WINDOW
	//userStakedSpaceBalance := sb.StakedSpaceBalance
	//allStakedSpaceBalance := getallStakedSpaceBalance()

	//maxUserInWindow := network_capacityInWindow * userStakedSpaceBalance / allStakedSpaceBalance
	//
	//lastUsage:=role.GetResourceUsage(sender)

//
//	long
//	oldNetUsage = accountCapsule.getNetUsage();
//	long
//	latestConsumeTime = accountCapsule.getLatestConsumeTime();
//	accountCapsule.setNetUsage(increase(oldNetUsage, 0, latestConsumeTime, now));
//	long
//	oldFreeNetUsage = accountCapsule.getFreeNetUsage();
//	long
//	latestConsumeFreeTime = accountCapsule.getLatestConsumeFreeTime();
//	accountCapsule.setFreeNetUsage(increase(oldFreeNetUsage, 0, latestConsumeFreeTime, now));
//	Map < String, Long > assetMap = accountCapsule.getAssetMap();
//	assetMap.forEach((assetName, balance)- >
//	{
//		long
//		oldFreeAssetNetUsage = accountCapsule.getFreeAssetNetUsage(assetName);
//		long
//		latestAssetOperationTime = accountCapsule.getLatestAssetOperationTime(assetName);
//		accountCapsule.putFreeAssetNetUsage(assetName,
//			increase(oldFreeAssetNetUsage, 0, latestAssetOperationTime, now));
//	});
//}
}

func getallStakedSpaceBalancee() uint64 {
	totalStakedBal := uint64(10000000000000)
	return totalStakedBal
}
