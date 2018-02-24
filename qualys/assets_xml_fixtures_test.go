package qualys

const (
	assetGroupsXMLSingleGroup = `
		<?xml version="1.0" encoding="UTF-8" ?>
		<!DOCTYPE ASSET_GROUP_LIST_OUTPUT SYSTEM "https://qualysapi.qualys.com/api/2.0/fo/asset/group/asset_group_list_output.dtd">
		<ASSET_GROUP_LIST_OUTPUT>
			<RESPONSE>
				<DATETIME>2016-10-05T19:00:22Z</DATETIME>
				<ASSET_GROUP_LIST>
					<ASSET_GROUP>
						<ID>1759735</ID>
						<TITLE><![CDATA[AG - Elastic Cloud Dynamic Perimeter]]></TITLE>
						<IP_SET>
							<IP>10.1.1.1</IP>
							<IP>10.10.10.11</IP>
						</IP_SET>
					</ASSET_GROUP>
				</ASSET_GROUP_LIST>
			</RESPONSE>
		</ASSET_GROUP_LIST_OUTPUT>
		<!-- CONFIDENTIAL AND PROPRIETARY INFORMATION. Qualys provides the QualysGuard Service "As Is," without any warranty of any kind. Qualys makes no warranty that the information contained in this report is complete or error-free. Copyright 2016, Qualys, Inc. //-->
	`

	assetGroupsXMLMultiGroups = `
		<?xml version="1.0" encoding="UTF-8" ?>
		<!DOCTYPE ASSET_GROUP_LIST_OUTPUT SYSTEM "https://qualysapi.qualys.com/api/2.0/fo/asset/group/asset_group_list_output.dtd">
		<ASSET_GROUP_LIST_OUTPUT>
			<RESPONSE>
				<DATETIME>2016-10-05T19:00:22Z</DATETIME>
				<ASSET_GROUP_LIST>
					<ASSET_GROUP>
						<ID>1759734</ID>
						<TITLE><![CDATA[AG - New]]></TITLE>
						<DEFAULT_APPLIANCE_ID>105102</DEFAULT_APPLIANCE_ID>
						<APPLIANCE_IDS>105102</APPLIANCE_IDS>
					</ASSET_GROUP>
					<ASSET_GROUP>
						<ID>1759735</ID>
						<TITLE><![CDATA[AG - Elastic Cloud Dynamic Perimeter]]></TITLE>
						<IP_SET>
							<IP_RANGE>10.10.10.3-10.10.10.6</IP_RANGE>
							<IP>10.10.10.14</IP>
						</IP_SET>
					</ASSET_GROUP>
				</ASSET_GROUP_LIST>
			</RESPONSE>
		</ASSET_GROUP_LIST_OUTPUT>
		<!-- CONFIDENTIAL AND PROPRIETARY INFORMATION. Qualys provides the QualysGuard Service "As Is," without any warranty of any kind. Qualys makes no warranty that the information contained in this report is complete or error-free. Copyright 2016, Qualys, Inc. //-->
	`

	assetGroupsAddIPsResponse = `
		<?xml version="1.0" encoding="UTF-8" ?>
		<!DOCTYPE SIMPLE_RETURN SYSTEM "https://qualysapi.qualys.com/api/2.0/simple_return.dtd">
		<SIMPLE_RETURN>
			<RESPONSE>
				<DATETIME>2016-10-12T14:16:22Z</DATETIME>
				<TEXT>Asset Group Updated Successfully</TEXT>
				<ITEM_LIST>
					<ITEM>
						<KEY>ID</KEY>
						<VALUE>1759735</VALUE>
					</ITEM>
				</ITEM_LIST>
			</RESPONSE>
		</SIMPLE_RETURN>
		`
)
