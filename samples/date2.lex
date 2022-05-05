<?xml version="1.0" encoding="UTF-8" ?>
<lexer xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
     xsi:noNamespaceSchemaLocation="lexer.xsd" 
>
    <name>Date1</name>
    <author>Guillaume de Lestanville</author>
    <description>Premier exemple simple d'analyse de date</description>
    <dateCreation>2020-12-10</dateCreation>
    <initial>0</initial>
    <rules>
        <rule id="1" from="0" to="0" repeat="3">
            <test>
                <charset>DIGIT</charset>
            </test>
            <concat>true</concat>
        </rule>
        <rule id="2" from="0" to="1">
            <test>
                <in>-/.</in>
            </test>
            <reset>true</reset>
            <action>Separator</action>
        </rule>
        <rule id="3" from="1" to="1" repeat="3">
            <test>
                <charset>DIGIT</charset>
            </test>
            <concat>true</concat>
        </rule>
        <rule id="4" from="1" to="2">
            <test>
                <in>-/.</in>
            </test>
            <reset>true</reset>
            <action>Separator</action>
        </rule>
        <rule id="5" from="2" to="2" repeat="3">
            <test>
                <charset>DIGIT</charset>
            </test>
            <concat>true</concat>
        </rule>
        <rule id="6" from="2" to="2">
            <test>
                <eos />
             </test>
             <final>true</final>
            <action>Separator</action>
        </rule>
    </rules>
</lexer>
