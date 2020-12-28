package common

import "testing"

func Test_Gzip(t *testing.T) {

	raw := []byte(`{"Key":"SF1031131137215","CMD":"waybill_op","LogID":"1600250136594297","Data":"{\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"waybillNo\":\"SF1031131137215\",\"sourceZoneCode\":\"573ND\",\"destZoneCode\":\"751\",\"meterageWeightQty\":1.0,\"realWeightQty\":0.5,\"quantity\":1.0,\"consigneeEmpCode\":\"92037265\",\"consignedTm\":\"2020-09-16 17:54:43\",\"cargoTypeCode\":\"T6\",\"limitTypeCode\":\"T6\",\"distanceTypeCode\":\"R10102\",\"transportTypeCode\":\"TR2\",\"expressTypeCode\":\"B1\",\"volume\":0.0,\"billLong\":0.0,\"billWidth\":0.0,\"billHigh\":0.0,\"versionNo\":5,\"lockVersionNo\":5,\"unitWeight\":\"kg\",\"productCode\":\"SE0004\",\"orderNo\":\"01901091611200879305838052\",\"updateTm\":\"2020-09-16 17:55:36\",\"createTm\":\"2020-09-16 11:20:11\",\"extJson\":\"{\\\"updateWaybillBarCode\\\":\\\"50\\\",\\\"autoConfirmWaybillFlag\\\":\\\"1\\\",\\\"orderSysSource\\\":\\\"BSP\\\",\\\"orderModifyTime\\\":\\\"2020-09-16 11:19:59:062\\\",\\\"orderCreateTrig\\\":-1,\\\"bar50ScanTm\\\":\\\"2020-09-16 17:54:43\\\",\\\"autoLoding\\\":\\\"3\\\",\\\"payFlg\\\":\\\"1\\\",\\\"monthlyCard\\\":\\\"5732469357\\\",\\\"bar50SelfSendFlag\\\":0,\\\"bar50Status\\\":0,\\\"swsVersionNo\\\":2,\\\"securityCodeTimestamp\\\":0,\\\"from50Weight\\\":true,\\\"from209Weight\\\":false,\\\"from50SourceWaybillNo\\\":false,\\\"barFillSourceZone\\\":true,\\\"barFillDestZone\\\":false,\\\"barFillWeight\\\":true,\\\"barFillBusinessType\\\":true,\\\"barFillDistanceType\\\":false,\\\"barFillConsigneeEmp\\\":true,\\\"barFillConsignedTm\\\":true,\\\"barFillPhone\\\":false,\\\"abnormalRedirectOrBackOrAddr\\\":false,\\\"usedBillFill\\\":true,\\\"awsmDeliveryTime\\\":0,\\\"waybillSendTopicFlag\\\":1,\\\"cabinetLimitType\\\":0,\\\"sgsOnBoxUpdate\\\":false,\\\"icsmCustomsModifyTm\\\":0,\\\"sopUpdatedProduct\\\":false,\\\"cxModifyTm\\\":0,\\\"from50StandardWeight\\\":false,\\\"bar50FillConsign\\\":true,\\\"updatePayMethodBy50Bar\\\":true,\\\"bar50FillCustoms\\\":false,\\\"bar50FillAddresseeAddr\\\":false,\\\"bar50FillVolume\\\":false,\\\"bar50ConsValueAndCode\\\":false,\\\"barFillLWH\\\":false,\\\"underCallOrder\\\":false,\\\"ext046047UpdatedByOrder\\\":false}\",\"actionJson\":\"{\\\"infoUpdate\\\":true}\",\"updateSource\":\"{\\\"pickup\\\":true,\\\"pickupTime\\\":\\\"2020-09-16 17:55:09\\\",\\\"sssAddr\\\":false,\\\"elecSign\\\":false,\\\"wemWaybill\\\":true,\\\"wemWaybillTime\\\":\\\"2020-09-16 17:55:14\\\",\\\"cxModify\\\":false,\\\"gisInfo\\\":true,\\\"gisInfoTime\\\":\\\"2020-09-16 11:20:16\\\",\\\"taxbillInfo\\\":false,\\\"awsm\\\":false,\\\"redirectOrBack\\\":false,\\\"uploadWaybill\\\":false,\\\"sws\\\":true,\\\"swsTime\\\":\\\"2020-09-16 17:55:36\\\",\\\"wbepWaybill\\\":false,\\\"orderInfo\\\":true,\\\"orderInfoTime\\\":\\\"2020-09-16 11:20:11\\\",\\\"pwmPackage\\\":false,\\\"esgCccp\\\":false,\\\"cmspAddr\\\":false,\\\"gisAssRds\\\":false,\\\"nonSwsRedInk\\\":false,\\\"eppUpdate\\\":false,\\\"bepUpdate\\\":false,\\\"bmrsUpdate\\\":false,\\\"cageUpdate\\\":false,\\\"pickupAgentNo\\\":false,\\\"swsAckBill\\\":false,\\\"pickupPackage\\\":false,\\\"sgsOnBox\\\":false,\\\"sopWaybill\\\":false,\\\"icsmCustoms\\\":false}\",\"genOrderFlag\":false,\"operationWaybillAddrCons\":{\"contactsId\":\"AAABdJTtRNVnnLtYdYhFhLstiJV3BKvN\",\"consignorCompName\":\"\",\"consignorAddr\":\"DE##EwA9TjC8qblwGIxpVHxGQa2ruUB8GvG2fkOY5euK8eZZHEuS\",\"consignorPhone\":\"DEEQAVTr46MHnUdPe4%2Bk1VW94NYxE%3D\",\"consignorContName\":\"徐爱光\",\"consignorMobile\":\"DEEQAVTr46MHnUdPe4%2Bk1VW94NYxE%3D\",\"addresseeCompName\":\"\",\"addresseeAddr\":\"DE##EwA9Tj1gpTLF5Bsk%2FGRAVE31Uu9HsQXXYo%2BnT6eLYjIerD6qqeLfg60xEuV7YBqy3T9SbWMce94zgA3CLSAjxb9VEkCVBUHuNSc70SgPLgCagIrwBZzMr6ViXf9PvWJDf4veCdqGWOsffE%2FkMMhmo0mR%2FPs%3D\",\"addresseePhone\":\"DEEQAVTt5YLPElJ1K2IrUmJHPs78Q%3D\",\"addresseeContName\":\"曾新华\",\"addresseeMobile\":\"DEEQAVTt5YLPElJ1K2IrUmJHPs78Q%3D\",\"consignorAddrNative\":\"5732469357\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:06\",\"createTm\":\"2020-09-16 11:20:11\",\"consignorCityCode\":\"573\",\"consignorDeptCode\":\"573ND\",\"addresseeCityCode\":\"751\",\"addresseeDeptCode\":\"751AB\",\"addresseeTransitCode\":\"751VA\",\"addresseeTeamCode\":\"751ABCL018\",\"consignorAreaCode\":\"571Y\",\"consignorHqCode\":\"CN02\",\"addresseeAreaCode\":\"020Y\",\"addresseeHqCode\":\"CN01\",\"senderProvince\":\"浙江省\",\"senderCity\":\"嘉兴市\",\"senderArea\":\"桐乡市\",\"senderAddr\":\"DE##EwA9TjC8qblwGIxpVHxGQa2ruUB8GvG2fkOY5euK8eZZHEuS\",\"receiverProvince\":\"广东省\",\"receiverCity\":\"韶关市\",\"receiverArea\":\"武江区\",\"receiverAddr\":\"DE##EwA9Tr%2FU6UMSQZlb16ewOu37LDTOxtGKTtj0v9dlPipRo3859Xdgb9nnJTOAuzj63%2BHIfWGXCnXeQdCAP3hCkg76qs2VM37TVe71joqTWQmP3SqSZ8oqj3NxJ9lR5cabuZz87Q%3D%3D\",\"addresseeAoiCode\":\"751ABCL029\",\"addresseeAoiType\":\"130000\",\"addresseeAoiId\":\"9326DB494EE64E618262A42D4EE48565\",\"consignorCountryCode\":\"CN\",\"addresseeCountryCode\":\"CN\",\"addresseeKeyWord\":\"朝阳村村民委员会\",\"covid19ConfirmedNum\":0,\"covid19SuspectedNum\":0,\"covid19Score\":0},\"operationWaybillAddrConsList\":[{\"contactsId\":\"AAABdJTtRNVnnLtYdYhFhLstiJV3BKvN\",\"consignorCompName\":\"\",\"consignorAddr\":\"DE##EwA9TjC8qblwGIxpVHxGQa2ruUB8GvG2fkOY5euK8eZZHEuS\",\"consignorPhone\":\"DEEQAVTr46MHnUdPe4%2Bk1VW94NYxE%3D\",\"consignorContName\":\"徐爱光\",\"consignorMobile\":\"DEEQAVTr46MHnUdPe4%2Bk1VW94NYxE%3D\",\"addresseeCompName\":\"\",\"addresseeAddr\":\"DE##EwA9Tj1gpTLF5Bsk%2FGRAVE31Uu9HsQXXYo%2BnT6eLYjIerD6qqeLfg60xEuV7YBqy3T9SbWMce94zgA3CLSAjxb9VEkCVBUHuNSc70SgPLgCagIrwBZzMr6ViXf9PvWJDf4veCdqGWOsffE%2FkMMhmo0mR%2FPs%3D\",\"addresseePhone\":\"DEEQAVTt5YLPElJ1K2IrUmJHPs78Q%3D\",\"addresseeContName\":\"曾新华\",\"addresseeMobile\":\"DEEQAVTt5YLPElJ1K2IrUmJHPs78Q%3D\",\"consignorAddrNative\":\"5732469357\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:06\",\"createTm\":\"2020-09-16 11:20:11\",\"consignorCityCode\":\"573\",\"consignorDeptCode\":\"573ND\",\"addresseeCityCode\":\"751\",\"addresseeDeptCode\":\"751AB\",\"addresseeTransitCode\":\"751VA\",\"addresseeTeamCode\":\"751ABCL018\",\"consignorAreaCode\":\"571Y\",\"consignorHqCode\":\"CN02\",\"addresseeAreaCode\":\"020Y\",\"addresseeHqCode\":\"CN01\",\"senderProvince\":\"浙江省\",\"senderCity\":\"嘉兴市\",\"senderArea\":\"桐乡市\",\"senderAddr\":\"DE##EwA9TjC8qblwGIxpVHxGQa2ruUB8GvG2fkOY5euK8eZZHEuS\",\"receiverProvince\":\"广东省\",\"receiverCity\":\"韶关市\",\"receiverArea\":\"武江区\",\"receiverAddr\":\"DE##EwA9Tr%2FU6UMSQZlb16ewOu37LDTOxtGKTtj0v9dlPipRo3859Xdgb9nnJTOAuzj63%2BHIfWGXCnXeQdCAP3hCkg76qs2VM37TVe71joqTWQmP3SqSZ8oqj3NxJ9lR5cabuZz87Q%3D%3D\",\"addresseeAoiCode\":\"751ABCL029\",\"addresseeAoiType\":\"130000\",\"addresseeAoiId\":\"9326DB494EE64E618262A42D4EE48565\",\"consignorCountryCode\":\"CN\",\"addresseeCountryCode\":\"CN\",\"addresseeKeyWord\":\"朝阳村村民委员会\",\"covid19ConfirmedNum\":0,\"covid19SuspectedNum\":0,\"covid19Score\":0}],\"operationWaybillMarkList\":[{\"labellingId\":\"AAABdJTtRNlBvHku7EJIFqMP7Mmra6m+\",\"labellingPattern\":\"DETAIL_FLG\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 11:20:11\",\"createTm\":\"2020-09-16 11:20:11\"},{\"labellingId\":\"AAABdJZW3fNO4B/fOydLCa2zTmmJR1xJ\",\"labellingPattern\":\"HAS_SERVICE_PROD_FLG\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:09\",\"createTm\":\"2020-09-16 17:55:09\"},{\"labellingId\":\"AAABdJZW8ydpDUyrINtLtKGfzX34WAmQ\",\"labellingPattern\":\"IS_CAL\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:14\",\"createTm\":\"2020-09-16 17:55:14\"}],\"operationWaybillFeeList\":[{\"feeId\":\"AAABdJZW8yaNhfFwGpNDdomsYXIDcZw+\",\"feeTypeCode\":\"4\",\"feeAmt\":7.47,\"gatherZoneCode\":\"573ND\",\"paymentTypeCode\":\"1\",\"settlementTypeCode\":\"2\",\"paymentChangeTypeCode\":\"0\",\"customerAcctCode\":\"5732469357\",\"currencyCode\":\"CNY\",\"serviceId\":\"AAABdJTtRNVrNaqgubBJzJAkrp+5T2tb\",\"gatherEmpCode\":\"92037265\",\"bizOwnerZoneCode\":\"573ND\",\"feeAmtInd\":7.47,\"feeIndType\":\"1\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:14\",\"createTm\":\"2020-09-16 17:55:14\"},{\"feeId\":\"AAABdJTtRNWPi4CEPBtCO53mRrRzNWRn\",\"feeTypeCode\":\"1\",\"feeAmt\":18.0,\"gatherZoneCode\":\"573ND\",\"paymentTypeCode\":\"1\",\"settlementTypeCode\":\"2\",\"paymentChangeTypeCode\":\"0\",\"customerAcctCode\":\"5732469357\",\"currencyCode\":\"CNY\",\"gatherEmpCode\":\"92037265\",\"bizOwnerZoneCode\":\"573ND\",\"feeAmtInd\":18.0,\"feeIndType\":\"1\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:14\",\"createTm\":\"2020-09-16 11:20:11\"},{\"feeId\":\"AAABdJZW8ydS4aSualFOZ6M6XTnlPIpn\",\"feeTypeCode\":\"5\",\"feeAmt\":498.0,\"gatherZoneCode\":\"573ND\",\"paymentTypeCode\":\"1\",\"settlementTypeCode\":\"2\",\"paymentChangeTypeCode\":\"0\",\"customerAcctCode\":\"5732469357\",\"currencyCode\":\"CNY\",\"serviceId\":\"AAABdJTtRNVrNaqgubBJzJAkrp+5T2tb\",\"gatherEmpCode\":\"92037265\",\"bizOwnerZoneCode\":\"573ND\",\"sourceCodeFeeAmt\":498.0,\"exchangeRate\":1.0,\"feeAmtInd\":498.0,\"feeIndType\":\"0\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:14\",\"createTm\":\"2020-09-16 17:55:14\"}],\"operationWaybillAdditionList\":[{\"extJson\":\"4:您好！受暴雨天气影响，部分机场及高速封闭，快件时效可能将有所增加，带来不便请您见谅\",\"additionalId\":\"AAABdJTtRNV9/sRIRcdIhInDVwTDsAXP\",\"additionalKey\":\"CONTROL_STRATEGY\",\"additionalValues\":\"4\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 11:20:11\",\"createTm\":\"2020-09-16 11:20:11\"},{\"additionalId\":\"AAABdJTtRNV1aFJfu9pMnJbavEvaieO3\",\"additionalKey\":\"HAS_OWN_WRAPPER_FLAG\",\"additionalValues\":\"0\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 11:20:11\",\"createTm\":\"2020-09-16 11:20:11\"},{\"additionalId\":\"AAABdJTtRNUEhzk1poxM7KmnBt+XYTsx\",\"additionalKey\":\"ORIGINAL_NUMBER\",\"additionalValues\":\"XS200916755749745307758612\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:09\",\"createTm\":\"2020-09-16 11:20:11\"},{\"additionalId\":\"AAABdJTtRNkdN3ldI2JFeZZEyihcMwWb\",\"additionalKey\":\"ORDER_SYS_SOURCE\",\"additionalValues\":\"BSP\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 11:20:11\",\"createTm\":\"2020-09-16 11:20:11\"},{\"additionalId\":\"AAABdJTtRNlN4qNCn29EgJNCDkRv6dc8\",\"additionalKey\":\"ORDER_TYPE\",\"additionalValues\":\"9\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:09\",\"createTm\":\"2020-09-16 11:20:11\"},{\"additionalId\":\"AAABdJTtRNltBcdITNVLwoYMTxuyqV+A\",\"additionalKey\":\"WAYBILL_NO_TYPE\",\"additionalValues\":\"WAYBILL_01_000\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 11:20:11\",\"createTm\":\"2020-09-16 11:20:11\"},{\"additionalId\":\"AAABdJZW3dNrepo6zUxBlqemjlTqWfOC\",\"additionalKey\":\"CONSIGNEE_PRE_MARK99_EMP_CODE\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:09\",\"createTm\":\"2020-09-16 11:20:11\"},{\"additionalId\":\"AAABdJZW3dP0EOFUzvFIdp2mMV7XBr4W\",\"additionalKey\":\"IS_SENSITIVE\",\"additionalValues\":\"0\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:09\",\"createTm\":\"2020-09-16 11:20:11\"},{\"additionalId\":\"AAABdJZW8yeLaG004x9AhLfFUFpxI2G/\",\"additionalKey\":\"IS_BILLING\",\"additionalValues\":\"Y\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:14\",\"createTm\":\"2020-09-16 17:55:14\"}],\"operationWaybillCustoms\":{\"exportId\":\"AAABdJTtRNWF5ZylhrRHZ6e4mf1B3aN4\",\"isUseUpstreamInvoice\":false,\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:06\",\"extJson\":\"{}\",\"createTm\":\"2020-09-16 11:20:11\"},\"operationWaybillCustomsList\":[{\"exportId\":\"AAABdJTtRNWF5ZylhrRHZ6e4mf1B3aN4\",\"isUseUpstreamInvoice\":false,\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:06\",\"extJson\":\"{}\",\"createTm\":\"2020-09-16 11:20:11\"}],\"operationWaybillServiceList\":[{\"serviceId\":\"AAABdJTtRNVrNaqgubBJzJAkrp+5T2tb\",\"serviceProdCode\":\"IN01\",\"attribute1\":\"498.00\",\"attribute2\":\"5732469357\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 11:20:11\",\"extJson\":\"{\\\"bar50Create\\\":false}\",\"createTm\":\"2020-09-16 11:20:11\"}],\"operationWaybillPackageList\":[{\"packageId\":\"AAABdJTtRNVxfUeZ2GdPSYBrHGVJwkIF\",\"packageNo\":\"SF1031131137215\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 11:20:11\",\"createTm\":\"2020-09-16 11:20:11\",\"operationWaybillConsignList\":[{\"consId\":\"AAABdJZW3dN7biqlFUhMF6gNDQyia/S0\",\"consName\":\"【皮尔卡丹】 高档男士休闲套装，货到付款！免费试穿  (颜色\",\"consQty\":\"1\",\"consValue\":0.0,\"packageId\":\"AAABdJTtRNVxfUeZ2GdPSYBrHGVJwkIF\",\"packageNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"updateTm\":\"2020-09-16 17:55:09\",\"extJson\":\"{}\",\"createTm\":\"2020-09-16 17:55:09\"}]}],\"operationTaxbillInfoList\":[],\"optWaybillAdditionExtList\":[{\"extId\":\"AAABdJTtRNW5FqRu72xG1pvm0ak9OaEs\",\"waybillId\":\"AAABdJTtRNXUvfbXMpVMI6omQdCaqZt3\",\"waybillNo\":\"SF1031131137215\",\"orderNo\":\"01901091611200879305838052\",\"attr021\":\"XS200916755749745307758612\",\"attr022\":\"5732469357\",\"createTime\":\"2020-09-16 11:20:11\",\"updateTm\":\"2020-09-16 11:20:11\"}],\"optWaybillSpecialHandlerList\":[],\"clientCode\":\"xwcp\",\"currentSource\":\"SWS\",\"deliveredType\":\"0\",\"barOpCode\":\"50\",\"consignTag\":\"C6\",\"limitTag\":\"T6\",\"receiptModelTag\":\"P3\",\"sortingModelTag\":\"S2\",\"containerModelTag\":\"W4\",\"deliveryModelTag\":\"D3\",\"productDisplayCode\":\"ProductDisplay1\",\"waybillLabelCode\":\"WayTab1\",\"standardMeterageWeight\":1.0,\"personalizeFlipRatio\":12000.0}","TimeStamp":"2020-09-16T17:55:36.61081322+08:00","SenderIP":"10.220.21.72"}`)

	r1, err := GzipMarshal(raw)
	if err != nil {
		t.Errorf("err=%s", err.Error())
		t.SkipNow()
	}

	r2, err := GzipUnmarshal(r1)
	if err != nil {
		t.Errorf("err=%s", err.Error())
		t.SkipNow()
	}

	if string(raw) != string(r2) {
		t.Error("should not happen")
		t.SkipNow()
	}

}