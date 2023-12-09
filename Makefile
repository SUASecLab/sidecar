RULES_IN := rules/rules.m4
RULES_OUT := rules/rules.json

$(RULES_OUT) : $(RULES_IN)
	m4 $(RULES_IN) > $(RULES_OUT)

.PHONY : clean
clean :
	-rm $(RULES_OUT)
