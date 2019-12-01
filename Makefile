include rules.mk
# for cloud run
APP = saas-one-func
BUILD_GO_DIR = $(shell pwd)/.tmp
BUILD_GO_MOD = $(BUILD_GO_DIR)/go.mod
BUILD_UI_DIR = $(shell pwd)/.ui
JUSTICE_UI = $(BUILD_UI_DIR)/justice-ui/dist/justiceUI


########################################
# Docker rules
########################################
.PHONY: docker
DOCKER_TAG = docker_$(APP)
docker: copy_repos_to_local ui 
	echo $(PWD)
	docker build --tag $(DOCKER_TAG) $(BUILD_GO_DIR)
	docker save $(DOCKER_TAG) -o $(DOCKER_TAG)
	docker image prune -f

########################################
# common Go rules
########################################

# copy local go repos to a temp directory
copy_repos_to_local: $(BUILD_GO_DIR) $(JUSTICE_UI)
	find . -maxdepth 1 -type f -exec cp {} $(BUILD_GO_DIR) \;
	cd $(BUILD_GO_DIR) && git clone git@github.com:extremenetworks/saas-xcentral.git
	# use this line instead when the xcentral repo is local
	#cp -r ../saas-xcentral/ $(BUILD_GO_DIR)
	rm -f $(BUILD_GO_DIR)/go.sum
	echo "module main" > ${BUILD_GO_MOD}
	echo "replace github.com/extremenetworks/saas-xcentral => ./saas-xcentral" >> ${BUILD_GO_MOD}
	rm -rf $(BUILD_GO_DIR)/*/.git
	find $(BUILD_GO_DIR) -depth -type d -name testdata\* -exec rm -rf {} \;
	cp -r $(JUSTICE_UI) $(BUILD_GO_DIR)

.PHONY: $(BUILD_GO_DIR) 
$(BUILD_GO_DIR):
	$(RM) -r $(BUILD_GO_DIR)
	mkdir -p $(BUILD_GO_DIR)

########################################
# UI rules
########################################
ui: $(JUSTICE_UI)
$(JUSTICE_UI): $(BUILD_UI_DIR)
	$(MAKE) -C  $(BUILD_UI_DIR)/justice-ui build

$(BUILD_UI_DIR):
	mkdir -p $(BUILD_UI_DIR)
	cd $(BUILD_UI_DIR) && git clone git@github.com:extremenetworks/justice-ui
	# use this line instead when the xcentral repo is local
	#cd $(BUILD_UI_DIR) && git clone $(BUILD_UI_DIR)/../../justice-ui 

clean:
	rm -rf $(BUILD_GO_DIR) $(BUILD_UI_DIR) $(APP) $(DOCKER_TAG)
