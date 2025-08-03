SRC_DIR := src
OBJ_DIR := obj
INC_DIR := include
BIN_DIR := bin

#CXX := g++
CXX := clang++
CXXFLAGS := -Wall -Wextra -pedantic -ansi -std=c++17 -stdlib=libc++ -O2 -I$(INC_DIR)

# Source files


ESMS_SRC := $(SRC_DIR)/esms.cpp $(SRC_DIR)/config.cpp $(SRC_DIR)/rosterplayer.cpp $(SRC_DIR)/comment.cpp \
			$(SRC_DIR)/penalty.cpp $(SRC_DIR)/report_event.cpp $(SRC_DIR)/esms.cpp $(SRC_DIR)/cond_utils.cpp \
			$(SRC_DIR)/teamsheet_reader.cpp $(SRC_DIR)/cond_action.cpp $(SRC_DIR)/cond_condition.cpp \
			$(SRC_DIR)/util.cpp $(SRC_DIR)/cond.cpp $(SRC_DIR)/mt.cpp $(SRC_DIR)/config.cpp \
			$(SRC_DIR)/tactics.cpp $(SRC_DIR)/anyoption.cpp $(SRC_DIR)/teamsheet.cpp
UPDTR_SRC := $(SRC_DIR)/rosterplayer.cpp $(SRC_DIR)/updtr.cpp $(SRC_DIR)/util.cpp \
			 $(SRC_DIR)/anyoption.cpp $(SRC_DIR)/config.cpp $(SRC_DIR)/comment.cpp $(SRC_DIR)/league_table.cpp
LGTABLE_SRC := $(SRC_DIR)/lgtable.cpp $(SRC_DIR)/league_table.cpp $(SRC_DIR)/util.cpp $(SRC_DIR)/anyoption.cpp
FIXTURES_SRC := $(SRC_DIR)/fixtures.cpp $(SRC_DIR)/util.cpp $(SRC_DIR)/anyoption.cpp
TSC_SRC := $(SRC_DIR)/tsc.cpp $(SRC_DIR)/rosterplayer.cpp $(SRC_DIR)/util.cpp $(SRC_DIR)/config.cpp \
		   $(SRC_DIR)/anyoption.cpp $(SRC_DIR)/teamsheet.cpp
ROSTER_CREATOR_SRC := $(SRC_DIR)/roster_creator.cpp $(SRC_DIR)/rosterplayer.cpp \
					  $(SRC_DIR)/anyoption.cpp $(SRC_DIR)/config.o $(SRC_DIR)/util.cpp

# Object files
ESMS_OBJ := $(ESMS_SRC:$(SRC_DIR)/%.cpp=$(OBJ_DIR)/%.o)
UPDTR_OBJ := $(UPDTR_SRC:$(SRC_DIR)/%.cpp=$(OBJ_DIR)/%.o)
LGTABLE_OBJ := $(LGTABLE_SRC:$(SRC_DIR)/%.cpp=$(OBJ_DIR)/%.o)
FIXTURES_OBJ := $(FIXTURES_SRC:$(SRC_DIR)/%.cpp=$(OBJ_DIR)/%.o)
TSC_OBJ := $(TSC_SRC:$(SRC_DIR)/%.cpp=$(OBJ_DIR)/%.o)
ROSTER_CREATOR_OBJ := $(ROSTER_CREATOR_SRC:$(SRC_DIR)/%.cpp=$(OBJ_DIR)/%.o)

# Targets
TARGETS := esms updtr lgtable fixtures tsc roster_creator

all: $(TARGETS)

esms: $(BIN_DIR)/esms
updtr: $(BIN_DIR)/updtr
lgtable: $(BIN_DIR)/lgtable
fixtures: $(BIN_DIR)/fixtures
tsc: $(BIN_DIR)/tsc
roster_creator: $(BIN_DIR)/roster_creator

$(BIN_DIR)/esms: $(ESMS_OBJ)
	$(CXX) $(CXXFLAGS) -o $@ $^

$(BIN_DIR)/updtr: $(UPDTR_OBJ)
	$(CXX) $(CXXFLAGS) -o $@ $^

$(BIN_DIR)/lgtable: $(LGTABLE_OBJ)
	$(CXX) $(CXXFLAGS) -o $@ $^

$(BIN_DIR)/fixtures: $(FIXTURES_OBJ)
	$(CXX) $(CXXFLAGS) -o $@ $^

$(BIN_DIR)/tsc: $(TSC_OBJ)
	$(CXX) $(CXXFLAGS) -o $@ $^

$(BIN_DIR)/roster_creator: $(ROSTER_CREATOR_OBJ)
	$(CXX) $(CXXFLAGS) -o $@ $^

# Compile .cpp to .o
$(OBJ_DIR)/%.o: $(SRC_DIR)/%.cpp
	@mkdir -p $(dir $@)
	$(CXX) $(CXXFLAGS) -c $< -o $@

# Clean
clean:
	rm -rf $(OBJ_DIR) $(TARGETS)

.PHONY: all clean esms updtr lgtable fixtures tsc roster_creator
